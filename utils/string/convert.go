package string

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

//字符串转Int
//    intStr：数字的字符串
func String2Int(intStr string) (intNum int) {
	intNum, _ = strconv.Atoi(intStr)
	return
}

//字符串转Int64
//    intStr：数字的字符串
func String2Int64(intStr string) (int64Num int64) {
	intNum, _ := strconv.Atoi(intStr)
	int64Num = int64(intNum)
	return
}

//字符串转Float64
//    floatStr：小数点数字的字符串
func String2Float64(floatStr string) (floatNum float64) {
	floatNum, _ = strconv.ParseFloat(floatStr, 64)
	return
}

//字符串转Float32
//    floatStr：小数点数字的字符串
func String2Float32(floatStr string) (floatNum float32) {
	floatNum64, _ := strconv.ParseFloat(floatStr, 32)
	floatNum = float32(floatNum64)
	return
}

//Int转字符串
//    intNum：数字字符串
func Int2String(intNum int) (intStr string) {
	intStr = strconv.Itoa(intNum)
	return
}

//Int64转字符串
//    intNum：数字字符串
func Int642String(intNum int64) (int64Str string) {
	//10, 代表10进制
	int64Str = strconv.FormatInt(intNum, 10)
	return
}

//Float64转字符串
//    floatNum：float64数字
//    prec：精度位数（不传则默认float数字精度）
func Float64ToString(floatNum float64, prec ...int) (floatStr string) {
	if len(prec) > 0 {
		floatStr = strconv.FormatFloat(floatNum, 'f', prec[0], 64)
		return
	}
	floatStr = strconv.FormatFloat(floatNum, 'f', -1, 64)
	return
}

//Float32转字符串
//    floatNum：float32数字
//    prec：精度位数（不传则默认float数字精度）
func Float32ToString(floatNum float32, prec ...int) (floatStr string) {
	if len(prec) > 0 {
		floatStr = strconv.FormatFloat(float64(floatNum), 'f', prec[0], 32)
		return
	}
	floatStr = strconv.FormatFloat(float64(floatNum), 'f', -1, 32)
	return
}

//二进制转10进制
func BinaryToDecimal(bit string) (num int) {
	fields := strings.Split(bit, "")
	lens := len(fields)
	var tempF float64 = 0
	for i := 0; i < lens; i++ {
		floatNum := String2Float64(fields[i])
		tempF += floatNum * math.Pow(2, float64(lens-i-1))
	}
	num = int(tempF)
	return
}

// GetString convert interface to string.
func GetString(v interface{}) string {
	switch result := v.(type) {
	case string:
		return result
	case []byte:
		return string(result)
	default:
		if v != nil {
			return fmt.Sprint(result)
		}
	}
	return ""
}

// GetInt convert interface to int.
func GetInt(v interface{}) int {
	switch result := v.(type) {
	case int:
		return result
	case int32:
		return int(result)
	case int64:
		return int(result)
	default:
		if d := GetString(v); d != "" {
			value, _ := strconv.Atoi(d)
			return value
		}
	}
	return 0
}

// GetInt64 convert interface to int64.
func GetInt64(v interface{}) int64 {
	switch result := v.(type) {
	case int:
		return int64(result)
	case int32:
		return int64(result)
	case int64:
		return result
	default:

		if d := GetString(v); d != "" {
			value, _ := strconv.ParseInt(d, 10, 64)
			return value
		}
	}
	return 0
}

// GetFloat64 convert interface to float64.
func GetFloat64(v interface{}) float64 {
	switch result := v.(type) {
	case float64:
		return result
	default:
		if d := GetString(v); d != "" {
			value, _ := strconv.ParseFloat(d, 64)
			return value
		}
	}
	return 0
}

// GetBool convert interface to bool.
func GetBool(v interface{}) bool {
	switch result := v.(type) {
	case bool:
		return result
	default:
		if d := GetString(v); d != "" {
			value, _ := strconv.ParseBool(d)
			return value
		}
	}
	return false
}