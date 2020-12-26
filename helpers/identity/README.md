# go-identity
用于解析身份证号码中的信息,年龄、性别、省份代码、城市代码和地区代码,根据地区代码获取地区名称


## 使用文档
```go
package main

import (
	"time"
	"fmt"
	"github.com/qq1060656096/go-identity"
)

func main() {
	id := identity.New("511521198305075558")
	err := id.Parse()
	// 解析错误
	if err != nil {
		fmt.Println(err)
	}
	// 获取省份代码
	id.GetProvince()
	// 获取城市代码
	id.GetCity()
	// 获取地区代码
	id.GetArea()
	// 获取生日
	id.GetBirthday()
	// 获取性别
	id.GetSex()
	if id.GetSex() == identity.SexMale {
		// 性别: 男
	} else {
		// 性别: 女
	}
	// 获取顺序码
	id.GetSequenceCode()
	// 获取验证码
	id.GetVerifyCode()
	
	// 18位、15位身份证判断
	switch identity.GetBitType() {
	case identity.BitType18:
		// 18位身份证处理
	case identity.BitType15:
		// 15位身份证处理
	}
	
	// 获取年龄
	AgeAt(id.GetBirthdayTime(), time.Date(2019, 3, 6, 0, 0, 0, 0, time.UTC))
	Age(id.GetBirthdayTime())
	// 检测出生日是不是闰年
	IsLeapYear(id.GetBirthdayTime())
}

```
