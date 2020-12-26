package base

import "fmt"

func StringList(args ...interface{}) []string {
	var x []string
	for _, arg := range args {
		switch arg := arg.(type) {
		case []string:
			x = append(x, arg...)
		case string:
			x = append(x, arg)
		default:
			panic("stringList: invalid argument of type " + fmt.Sprintf("%T", arg))
		}
	}
	return x
}
