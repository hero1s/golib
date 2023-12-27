package toolkit

import "reflect"

func IsFunction(f interface{}) bool {
	actual := reflect.TypeOf(f)
	return actual.Kind() == reflect.Func && actual.NumIn() == 0 && actual.NumOut() > 0
}

func IsChan(a interface{}) bool {
	if IsNil(a) {
		return false
	}
	return reflect.TypeOf(a).Kind() == reflect.Chan
}

func IsNil(a interface{}) bool {
	if a == nil {
		return true
	}

	switch reflect.TypeOf(a).Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return reflect.ValueOf(a).IsNil()
	}

	return false
}
