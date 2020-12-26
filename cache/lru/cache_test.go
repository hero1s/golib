package lru

import (
	"github.com/hero1s/golib/log"
	"testing"
)

func TestNewCache(t *testing.T) {
	cache := NewCache(2)

	cache.Put("1", "one")
	log.Infof("%v", cache.Get("1"))

	log.Warnf("===============")
	cache.Put("2", "two")
	log.Infof("%v", cache.Get("1"))

	log.Warnf("===============")
	cache.Put("3", "three")
	log.Infof("%v", cache.Get("2"))
	log.Infof("%v", cache.Get("3"))
	log.Infof("%v", cache.Get("3"))
	log.Infof("%v", cache.Get("1"))

	log.Warnf("===============")
	cache.Put("2", "two")
	log.Infof("%v", cache.Get("3"))
	log.Infof("%v", cache.Get("1"))
}
