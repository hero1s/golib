package retry

import (
	"errors"
	"github.com/hero1s/golib/log"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	err := Retry(func() error {
		log.Debugf("retry func")
		return errors.New("please retry")
	}, 3, 2*time.Second)
	if err != nil {
		log.Errorf(err.Error())
	}
}
