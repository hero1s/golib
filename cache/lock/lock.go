package lock

import (
	"context"
	"github.com/hero1s/golib/cache"
	"github.com/hero1s/golib/log"
	"github.com/bsm/redislock"
	"time"
)

func TryLock(lockName string, ttl, waitTime time.Duration) (*redislock.Lock, error) {
	locker := redislock.New(cache.Redis)
	ctx, cancel := context.WithTimeout(context.Background(), waitTime)
	defer cancel()
	opt := redislock.Options{
		RetryStrategy: redislock.ExponentialBackoff(100*time.Millisecond, waitTime),
		Context:       ctx,
	}
	lock, err := locker.Obtain(lockName, ttl, &opt)
	if err == redislock.ErrNotObtained {
		log.Errorf("Could not obtain lock! lockerName:%v", lockName)
	} else if err != nil {
		log.Errorf(err.Error())
	}
	return lock, err
}
