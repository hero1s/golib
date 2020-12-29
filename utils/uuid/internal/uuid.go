package internal

import (
	"errors"
	"fmt"
	"github.com/hero1s/golib/log"
	"github.com/go-redis/redis/v7"
	"io"
	"sync"
	"sync/atomic"
)

// 思想就是低位和高位都是跳动变化，这样达到不是一个完全顺序的数
// 低位满了之后，高位重新生成一位,直到64位都用完
//  ---高52位---|---低12位---
const (
	// RenewInterval indicates how often renew retries are performed
	// 低12位使用ID
	Renew12Interval uint64 = 0xDFF
	// 64位都用完了 15个F
	PanicValue uint64 = 0x7FFFFFFFFFFFFFFF
)

// UUID is for internal use only.
type UUID struct {
	sync.Mutex
	N           uint64
	Tag         string
	Renew       func() error
	H52Verifier func(h52 uint64) error // 更高位52位
}

type Option func(*UUID)

func NewUUID(tag string, opts ...Option) *UUID {
	u := &UUID{Tag: tag}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

// Next is for internal use only.
func (u *UUID) Next() uint64 {
	x := atomic.AddUint64(&u.N, 1)
	// 值已经大于64位的最大值了
	if x >= PanicValue {
		log.Errorf("<uuid> 已经达到极限值, tag: %s", u.Tag)
		panic(fmt.Errorf("<uuid> 已经达到极限值, tag: %s", u.Tag))
	}
	// 低12位已经满了
	if x&Renew12Interval >= Renew12Interval {
		err := u.RenewNow()
		if err != nil {
			log.Errorf("<uuid> renew 52 failed. tag: %s, reason: %+v", u.Tag, err)
		} else {
			log.Debugf("<uuid> renew 52 succeeded. tag: %s", u.Tag)
		}
	}
	return x
}

// RenewNow reacquires the high 52 bits from your data store immediately
func (u *UUID) RenewNow() error {
	u.Lock()
	renew := u.Renew
	u.Unlock()

	return renew()
}

// Reset is for internal use only.
func (u *UUID) Reset(n uint64) {
	atomic.StoreUint64(&u.N, n)
}

// VerifyH52 is for internal use only.
func (u *UUID) VerifyH52(h52 uint64) error {
	if h52 == 0 {
		return errors.New("the high 52 bits should not be 0. tag: " + u.Tag)
	}

	if h52 > 0x7FFFFFFFFFFFF {
		return errors.New("the high 52 bits should not exceed 0x0FFFFFFF. tag: " + u.Tag)
	}

	if u.H52Verifier != nil {
		if err := u.H52Verifier(h52); err != nil {
			return err
		}
	}
	return nil
}

type NewClient func() (client redis.Cmdable, autoDisconnect bool, err error)

func (u *UUID) LoadH52FromRedis(newClient NewClient, key string) error {
	if len(key) == 0 {
		return errors.New("key cannot be empty. tag: " + u.Tag)
	}

	client, autoDisconnect, err := newClient()
	if err != nil {
		return err
	}
	if autoDisconnect {
		defer func() {
			closer := client.(io.Closer)
			_ = closer.Close()
		}()
	}

	n, err := client.Incr(key).Result()
	if err != nil {
		return err
	}
	h52 := uint64(n)
	if err = u.VerifyH52(h52); err != nil {
		return err
	}

	u.Reset(h52 << 12)
	log.Infof("<uuid> new h52: %d. tag: %s\n", h52, u.Tag)

	u.Lock()
	defer u.Unlock()

	if u.Renew != nil {
		return nil
	}
	u.Renew = func() error {
		return u.LoadH52FromRedis(newClient, key)
	}
	return nil
}
