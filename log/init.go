package log

import (
	"fmt"
	"github.com/hero1s/golib/helpers/ip"
	"github.com/hero1s/golib/helpers/webhook"
	"github.com/hero1s/golib/log/conf"
	"github.com/hero1s/golib/log/plugins/zaplog"
	"os"
	"runtime"
	"time"
)

// 默认
var l Loger = zaplog.New()
var larkUrl = ""
var localIP = "127.0.0.1"
var curDir = ""

func CurLogger() *Loger {
	return &l
}

// 设置
func InitLogger(opts ...conf.Option) {
	l = zaplog.New(opts...)
}

func SetLarkUrl(url string) {
	larkUrl = url
	localIP = ip.InternalIP()
	curDir, _ = os.Getwd()
}

// 快捷使用
func EasyInit(projectName, filePath string, stdout bool) {
	InitLogger(conf.WithProjectName(projectName),
		conf.WithFilename(filePath),
		conf.WithLogType(conf.LogJsontype),
		conf.WithIsStdOut(stdout))
}

// 目前只有zap生效
func SetLogLevel(level conf.Level) {
	l.SetLogLevel(level)
}

// 日志同步写入
// 目前只有zap生效
func Sync() {
	l.Sync()
}

// 日志等级 调试时使用
func Debug(keysAndValues ...interface{}) {
	l.Debug(keysAndValues...)
}

// 日志等级 提示时使用
func Info(keysAndValues ...interface{}) {
	l.Info(keysAndValues...)
}

// 日志等级 警告时使用
func Warn(keysAndValues ...interface{}) {
	l.Warn(keysAndValues...)
}

// 日志等级 错误时使用
func Error(keysAndValues ...interface{}) {
	l.Error(keysAndValues...)
	SendLarkf("%v", keysAndValues)
}

// 日志等级 恐慌时使用
func Panic(keysAndValues ...interface{}) {
	l.Panic(keysAndValues...)
	SendLarkf("%v", keysAndValues)
}

// 日志等级 致命时使用
func Fatal(keysAndValues ...interface{}) {
	l.Fatal(keysAndValues...)
	SendLarkf("%v", keysAndValues)
}

// 日志等级 详细结构类型,调试利器
func Dump(keysAndValues ...interface{}) {
	l.Dump(keysAndValues...)
}

// 调试的
func Debugf(template string, args ...interface{}) {
	l.Debugf(template, args...)
}

// 提示的
func Infof(template string, args ...interface{}) {
	l.Infof(template, args...)
}

// 警告的
func Warnf(template string, args ...interface{}) {
	l.Warnf(template, args...)
}

// 错误的
func Errorf(template string, args ...interface{}) {
	l.Errorf(template, args...)
	SendLarkf(template, args...)
}

// 恐慌的
func Panicf(template string, args ...interface{}) {
	l.Panicf(template, args...)
	SendLarkf(template, args...)
}

// 致命的
func Fatalf(template string, args ...interface{}) {
	l.Fatalf(template, args...)
	SendLarkf(template, args...)
}

// 发送错误日志到飞书
func SendLarkf(template string, args ...interface{}) {
	if larkUrl != "" {
		// 获取调用文件名和行号
		_, file, line, ok := runtime.Caller(2) // 调用层级为1，表示直接调用此函数的地方
		if !ok {
			file = "unknown"
			line = 0
		}
		message := fmt.Sprintf(template, args...)
		go func(file string, line int, message string) {
			timestamp := time.Now().Format("2006-01-02 15:04:05") // 标准时间格式
			// 构造带时间戳、IP、路径、文件名及行号的日志信息
			messageWithMeta := fmt.Sprintf("[%s][%s][%s][%s:%d] %s", timestamp, localIP, curDir, file, line, message)
			webhook.SendLark(messageWithMeta, larkUrl)
		}(file, line, message)
	}
}
