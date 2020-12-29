package threading

import (
	"github.com/hero1s/golib/utils/rescue"
)

/*安全go执行函数，注意传参延迟的问题，
不要go函数中引用可变参数。特别是for循环中
尽量在函数内部GoSafe
*/
func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	defer rescue.Recover()

	fn()
}
