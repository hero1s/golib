package rescue

import (
	"github.com/hero1s/golib/log"
	"runtime/debug"
)

func Recover(pinicCallback func(stack string), cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}
	if p := recover(); p != nil {
		log.Errorf("call func panic occure,err:%v", p)
		stack := string(debug.Stack())
		log.Error("stack:%v", stack)
		if pinicCallback != nil {
			pinicCallback(stack)
		}
	}
}
