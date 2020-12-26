package regex

import (
	"github.com/hero1s/golib/log"
	"regexp"
)

func CheckPhone(areaCode string, phone string) bool {
	if areaCode == "+86" {
		return CheckChinaPhone(phone)
	}
	reg := `^[0-9\-]*$`
	return checkPhone(reg, phone)
}

func CheckChinaPhone(phone string) bool {
	reg := `^(13[0-9]|14[579]|15[0-3,5-9]|16[6]|17[0135678]|18[0-9]|19[89])\d{8}$`
	return checkPhone(reg, phone)
}

func CheckHKPhone(phone string) bool {
	reg := `^([5689])\d{7}$`
	return checkPhone(reg, phone)
}
func CheckTWPhone(phone string) bool {
	reg := `^(9|09)\d{8}$`
	return checkPhone(reg, phone)
}
func CheckMacaoPhone(phone string) bool {
	reg := `^(6)\d{8}$`
	return checkPhone(reg, phone)
}

func checkPhone(reg string, phone string) bool {
	res, err := regexp.MatchString(reg, phone)
	if err != nil {
		log.Debugf("regexp phone error,cause:%s", err)
		return false
	}
	return res
}
