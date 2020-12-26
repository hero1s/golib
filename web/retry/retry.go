package retry

import (
	"github.com/hero1s/golib/log"
	"time"
)

// Retry 重试 func 最大次数，间隔
func Retry(callback func() error, maxRetries int, interval time.Duration) error {
	var err error
	for i := 1; i <= maxRetries; i++ {
		if err = callback(); err != nil {
			log.Errorf("Retry(%d) error(%+v)", i, err)
			time.Sleep(interval)
			continue
		}
		return nil
	}
	return err
}
