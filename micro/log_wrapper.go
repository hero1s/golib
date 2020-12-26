package micro

import (
	"context"
	"fmt"
	"github.com/hero1s/golib/log"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
	"go.uber.org/zap"
	"runtime"
	"time"
)

func LogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		begin := time.Now()
		meta, _ := metadata.FromContext(ctx)
		defer func() {
			log.Info("access log",
				zap.Any("meta", meta),
				zap.String("ip", meta["Local"]+"--"+meta["Remote"]),
				zap.String("method", req.Method()),
				zap.String("path", meta["Micro-Service"]+"."+meta["Micro-Method"]),
				zap.String("Request-Id", req.Header()["Request-Id"]),
				zap.String("queries", req.Endpoint()),
				zap.Any("body", req.Body()),
				zap.Any("rsp", rsp),
				zap.Duration("duration", time.Now().Sub(begin)))
		}()
		return fn(ctx, req, rsp)
	}
}

func LogRecover(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		defer func() {
			if err := recover(); err != nil {
				var trace string
				for i := 1; ; i++ {
					if _, f, l, got := runtime.Caller(i); !got {
						break
					} else {
						trace += fmt.Sprintf("%s:%d;", f, l)
					}
				}

				log.Error(
					zap.String("recover exception.server", req.Endpoint()),
					zap.Any("reqBody", req.Body()),
					zap.Any("error", err),
					zap.String("trace", trace))
			}
		}()
		return fn(ctx, req, rsp)
	}
}

func LogClientWrap(c client.Client) client.Client {
	return &logClientWrapper{c}
}

type logClientWrapper struct {
	client.Client
}

func (l *logClientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	log.Infof("client[%s], method[%s], params[%+v]", req.Service(), req.Method(), req.Body())
	return l.Client.Call(ctx, req, rsp)
}
