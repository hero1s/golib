package micro

import (
	"github.com/hero1s/golib/log"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
	"time"
)

type JaegerConf struct {
	ServiceName string `json:"service_name"`
	Addr        string `json:"addr"`
}

// NewTracer 创建一个jaeger Tracer
func NewTracer(conf JaegerConf) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: conf.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	sender, err := jaeger.NewUDPTransport(conf.Addr, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)

	return tracer, closer, err
}

// 初始化全局tracer
func InitGlobalTracer(conf JaegerConf) {
	t, _, err := NewTracer(conf)
	if err != nil {
		log.Error(err.Error())
	}
	opentracing.SetGlobalTracer(t)
}
