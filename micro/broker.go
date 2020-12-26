package micro

import (
	"github.com/hero1s/golib/log"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-plugins/broker/kafka/v2"
	"github.com/micro/go-plugins/broker/mqtt/v2"
	"github.com/micro/go-plugins/broker/nsq/v2"
	"github.com/micro/go-plugins/broker/redis/v2"
	"strings"
)

type Broker interface {
	Address() []string
}

func nsqBroker(c *NsqBroker) broker.Broker {
	return nsq.NewBroker(func(options *broker.Options) {
		options.Addrs = c.Addrs
	})
}

func redisBroker(c *RedisBroker) broker.Broker {
	return redis.NewBroker(func(options *broker.Options) {
		options.Addrs = c.Addrs
	})
}

func mqttBroker(c Broker) broker.Broker {
	return mqtt.NewBroker(func(options *broker.Options) {
		options.Addrs = c.Address()
		options.Secure = true
	})
}

// 初始化 broker
func kafkaBroker(c Broker) broker.Broker {
	var kafkaBroker broker.Broker
	if len(c.Address()) < 1 || strings.TrimSpace(c.Address()[0]) == "" {
		log.Errorf("kafka address is error:%v",c)
		kafkaBroker = nil
	} else {
		kafkaBroker = kafka.NewBroker(func(options *broker.Options) {
			options.Addrs = c.Address()
		})
	}
	return kafkaBroker
}
