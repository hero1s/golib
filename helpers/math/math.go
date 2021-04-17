package math

import (
	"math"
	"math/rand"
	"time"
	"github.com/shopspring/decimal"
)

func init() {
	rand.Seed(time.Now().Unix())
}

//随机值在闭区间[min,max]
func Random(min, max int64) int64 {
	max += 1
	return rand.Int63n(max-min) + min
}
func RemoveDuplicateInt(list []int64) []int64 {
	var x []int64
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}
func RemoveDuplicateStr(list []string) []string {
	var x []string = []string{}
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

// int数组合并，去重复
func MergeSliceRemoveDuplicate(slice1, slice2 []int) (merged []int) {
	var dupMap = make(map[int]int)
	slice1 = append(slice1, slice2...)
	for _, v := range slice1 {
		length := len(dupMap)
		dupMap[v] = 1
		if len(dupMap) != length {
			merged = append(merged, v)
		}
	}
	return merged
}

func MinInt(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func MaxInt(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

//保留N位小数
func Decimal(value float64,places int32) float64 {
	v1, _ := decimal.NewFromFloat(value).Round(places).Float64()
	return v1
}

//打乱一个uint数组
func RandomUint(uints []uint64) {
	for i := len(uints) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		uints[i], uints[num] = uints[num], uints[i]
	}
}

//打乱一个string数组
func RandomString(strings []string) {
	for i := len(strings) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		strings[i], strings[num] = strings[num], strings[i]
	}
}

