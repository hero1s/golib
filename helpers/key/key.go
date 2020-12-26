package key

import (
	"math/rand"
	"strings"
)

// 一个自定义的秘钥或随机字符串生成器
const(
	StrSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	HexSource = "0123456789ABCDEF"
	NumSource = "0123456789"
)

// 用内置字符串生成随机字符串或秘钥
func RandomStr(length int) string {
	sb := [] string{}
	if length > 0 {
		for i :=0; i< length;i++  {
			sb = append(sb, string(StrSource[rand.Intn(len(StrSource))]))
		}
	}
	return strings.Join(sb,"")
}

// 生成十六进制随机串
func RandomHex(length int) string {
	sb := [] string{}
	if length > 0 {
		for i :=0; i< length;i++  {
			sb = append(sb, string(HexSource[rand.Intn(len(HexSource))]))
		}
	}
	return strings.Join(sb,"")
}

// 生成纯数字的随机串
func RandomNum(length int) string {
	sb := [] string{}
	if length > 0 {
		for i :=0; i< length;i++  {
			sb = append(sb, string(NumSource[rand.Intn(len(NumSource))]))
		}
	}
	return strings.Join(sb,"")
}
