package log

import (
	"github.com/hero1s/golib/log/conf"
)

//定义接口
type Loger interface {
	//key value
	Debug(...interface{}) //调试的
	Info(...interface{})  //提示的
	Warn(...interface{})  //警告的
	Error(...interface{}) //错误的
	Panic(...interface{}) //恐慌的
	Fatal(...interface{}) //致命的
	//sprintf 不建议使用性能低，非json格式不方便接入第三方日志库
	Debugf(template string, args ...interface{}) //调试的
	Infof(template string, args ...interface{})  //提示的
	Warnf(template string, args ...interface{})  //警告的
	Errorf(template string, args ...interface{}) //错误的
	Panicf(template string, args ...interface{}) //恐慌的
	Fatalf(template string, args ...interface{}) //致命的
	Dump(...interface{})                         //详细结构类型,调试利器
	Sync()                                       //同步
	SetLogLevel(conf.Level)                      //可以随机设置日志级别的
	Write(p []byte) (n int, err error)           //io.Writer
}
