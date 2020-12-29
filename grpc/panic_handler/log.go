package panic_handler

import (
	"fmt"
	"github.com/hero1s/golib/log"
	"os"
	"runtime"
	"runtime/debug"
)

func LogPanicDump(r interface{}) {
	fmt.Fprintf(os.Stderr, string(debug.Stack()))
}

func LogPanicStackMultiLine(r interface{}) {
	callers := []string{}
	for i := 0; true; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		callers = append(callers, fmt.Sprintf("%d: %v:%v", i, file, line))
	}
	if len(callers) > 0 {
		log.Errorf("Recovered from panic: %#v (%v) in %s", r, r, callers[0])
	}
	log.Warnf("StackTrace:")
	for i := 0; len(callers) > i; i++ {
		log.Errorf("  %s", callers[i])
	}
}