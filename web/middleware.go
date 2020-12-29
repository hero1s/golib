package web

import (
	"bytes"
	"context"
	"github.com/hero1s/golib/helpers/token"
	"github.com/hero1s/golib/log"
	"github.com/hero1s/golib/utils/uuid"
	"github.com/hero1s/golib/web/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// CORS gin middleware cors
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") // 请求头部
		if origin == "" {
			origin = c.Request.Host
		}
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, X-CSRF-Token, authorization, sign, appid, ts")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

// Recovery gin middleware recovery
func Recovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					log.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					log.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					log.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				control := response.Control{}
				control.InternalError(c, "Server Error")
			}
		}()
		c.Next()
	}
}

// token 校验中间件
func JWTAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		m, err := token.DecodeToken(context.Request)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"status":  -1,
				"message": err.Error(),
			})
			context.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		context.Set("claims", m)
		context.Next()
	}
}

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.Request.Header.Get("Request-Id")
		if requestId == "" {
			requestId = uuid.NewUuid()
		}
		c.Set("Request-Id", requestId)
		c.Header("Request-Id", requestId)
		c.Next()
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// 重写, 将同样的数据写一份保存到 body 中
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

const contextTracerKey = "Tracer-context"

// sf sampling frequency
var sf = 100

// SetSamplingFrequency 设置采样频率
// 0 <= n <= 100
func SetSamplingFrequency(n int) {
	sf = n
}

// 日志记录到文件
func LogRequest(openTrace bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.EscapedPath()
		method := c.Request.Method
		ip := c.ClientIP()

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		// 读取后写回
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw
		start := time.Now()

		//是否上报openTrace
		var sp opentracing.Span
		var nsf int
		if openTrace {
			sp = opentracing.GlobalTracer().StartSpan(c.Request.URL.Path)
			tracer := opentracing.GlobalTracer()
			md := make(map[string]string)
			nsf = sf
			spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
			if err == nil {
				sp = opentracing.GlobalTracer().StartSpan(c.Request.URL.Path, opentracing.ChildOf(spanCtx))
				tracer = sp.Tracer()
				nsf = 100
			}
			defer sp.Finish()

			if err := tracer.Inject(sp.Context(),
				opentracing.TextMap,
				opentracing.TextMapCarrier(md)); err != nil {
				log.Error(err)
			}

			ctx := context.TODO()
			ctx = opentracing.ContextWithSpan(ctx, sp)
			ctx = metadata.NewContext(ctx, md)
			c.Set(contextTracerKey, ctx)
		}

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		code, message := c.Writer.Status(), blw.body.Bytes()

		log.Info(
			zap.Int("status", code),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", string(bodyBytes)),
			zap.String("ip", ip),
			zap.Any("head", c.Request.Header),
			zap.String("resp", string(message)),
			zap.Duration("latency", latency),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		)

		if openTrace {
			ext.HTTPStatusCode.Set(sp, uint16(code))
			ext.HTTPMethod.Set(sp, method)
			ext.HTTPUrl.Set(sp, path)
			ext.PeerAddress.Set(sp, ip)
			sp.SetTag("head", c.Request.Header)
			sp.SetTag("query", string(bodyBytes))
			sp.SetTag("resp", string(message))
			if code >= http.StatusInternalServerError {
				ext.Error.Set(sp, true)
			} else if rand.Intn(100) > nsf {
				ext.SamplingPriority.Set(sp, 0)
			}
		}
	}
}
