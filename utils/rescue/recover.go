package rescue

import (
	"git.moumentei.com/plat_go/golib/log"
	"runtime/debug"
)

type NotifyFunc func(stack string)

var notifyFunc NotifyFunc = nil

//设置报警函数(钉钉/邮件)
func SetNotifyFunc(f NotifyFunc) {
	notifyFunc = f
}

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}
	if p := recover(); p != nil {
		log.Errorf("call func panic occure,err:%v", p)
		log.Error("stack:%v", string(debug.Stack()))
		if notifyFunc != nil {
			notifyFunc(string(debug.Stack()))
		}
	}
}
