package micro

import (
	hystrix2 "github.com/afex/hystrix-go/hystrix"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix/v2"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	opentracing2 "github.com/opentracing/opentracing-go"
)

func InitClient(cliName, version string, registry *EtcdRegistry, broker *KafkaBroker, fn func(client.Client)) {
	var (
		s    micro.Service
		opts []micro.Option
	)
	opts = []micro.Option{
		micro.Name(cliName),
		//micro.WrapClient(LogClientWrap),
		micro.WrapClient(hystrix.NewClientWrapper()),
		micro.Version(version),
	}

	//处理命令行参数
	opts = append(opts, micro.Flags(&cli.StringFlag{
		Name:     "conf",
		Usage:    "config filename",
		Value:    "./conf.toml",
		Required: false,
	}),
		/*		micro.Action(func(c *cli.Context) error {
				c.String("conf")
				return nil
			}),*/)

	if registry != nil {
		reg := etcdRegistry(registry)
		s := selector.NewSelector(func(options *selector.Options) { //暂时写死轮询 toney
			options.Registry = reg
		})
		opts = append(opts, micro.Selector(s))
		opts = append(opts, micro.Registry(reg))
	}
	if broker != nil {
		opts = append(opts, micro.Broker(kafkaBroker(broker)))
	}

	if opentracing2.IsGlobalTracerRegistered() {
		opts = append(opts, micro.WrapClient(opentracing.NewClientWrapper(opentracing2.GlobalTracer())))
	}

	//全局限流
	hystrix2.DefaultMaxConcurrent = 50
	hystrix2.DefaultTimeout = 2000
	//单服务限流
	//hystrix2.ConfigureCommand("",hystrix2.CommandConfig{MaxConcurrentRequests: 10,Timeout: 1000})

	s = micro.NewService(opts...)
	//s.Init()
	fn(s.Client())
}
