package identity

import (
	"time"
)

// AgeAt 获取特定时间之间的年龄
func AgeAt(birthDate time.Time, nowDate time.Time) int {
	age := nowDate.Year() - birthDate.Year()
	if nowDate.Month() <= birthDate.Month() && nowDate.Day() < birthDate.Day(){
		age -= 1
	}
	return age
}

// Age 根据出生日期获取年龄
func Age(birthdayTime time.Time) int {
	return AgeAt(birthdayTime, time.Now())
}

// IsLeapYear 检测日期是否是闰年
func IsLeapYear(date time.Time) bool {
	y := date.Year()
	switch {
	case y % 400 == 0:
		return true
	case y % 100 == 0:
		return true
	case y % 4 == 0:
		return true
	}
	return false
}