package micro

import (
	"github.com/hero1s/golib/log"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"go.uber.org/zap"
	"strings"
)

type Registry interface {
	Address() []string
	UserName() string
	PassWord() string
}

func etcdRegistry(c *EtcdRegistry) (reg registry.Registry) {
	var Registry registry.Registry
	if len(c.Addrs) < 1 || strings.TrimSpace(c.Addrs[0]) == "" {
		log.Error(zap.Any("etcd error:", c))
		return nil
	}
	if strings.TrimSpace(c.User) != "" && strings.TrimSpace(c.Pass) != "" {
		Registry = etcd.NewRegistry(func(options *registry.Options) {
			options.Addrs = c.Addrs
			etcd.Auth(c.User, c.Pass)
		})
	} else {
		Registry = etcd.NewRegistry(func(options *registry.Options) {
			options.Addrs = c.Addrs
		})
	}
	return Registry
}
