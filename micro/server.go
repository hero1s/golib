package micro

import (
	"fmt"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	limiter "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	wrapperTrace "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	opentracing2 "github.com/opentracing/opentracing-go"
)

func InitServer(name, version, addr string, registry *EtcdRegistry, broker *KafkaBroker, fn func(s micro.Service)) {
	var (
		s    micro.Service
		opts []micro.Option
	)
	opts = []micro.Option{
		micro.Name(name),
		micro.WrapHandler(LogWrapper),
		micro.WrapHandler(LogRecover),
		micro.Version(version),
		micro.WrapHandler(limiter.NewHandlerWrapper(100)), //QPS = 100 限流保护服务器

		//过滤conf参数
		micro.Flags(&cli.StringFlag{
			Name:     "conf",
			Usage:    "config filename",
			Value:    "./conf.toml",
			Required: false,
		}),
	}
	if addr != "" {
		opts = append(opts, micro.Address(addr)) //这里可以指定一个端口 ":8080"，也可以cmd指定多个端口提供服务)
	}
	if registry != nil {
		opts = append(opts, micro.Registry(etcdRegistry(registry)))
	}
	if broker != nil {
		opts = append(opts, micro.Broker(kafkaBroker(broker)))
	}
	if opentracing2.IsGlobalTracerRegistered() {
		opts = append(opts, micro.WrapHandler(wrapperTrace.NewHandlerWrapper(opentracing2.GlobalTracer())))
	}

	s = micro.NewService(opts...)
	//s.Init()
	fn(s)

	go func() {
		if err := s.Run(); err != nil {
			panic(fmt.Sprintf("[%s] micro server run error(%+v).", name, err))
		}
	}()
}
