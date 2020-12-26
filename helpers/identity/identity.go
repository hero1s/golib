package identity

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

const (
	// 身份证位数(15位或者18位)
	BitType = "bitType"
	// 省份
	Province = "province"
	// 城市
	City = "city"
	// 地区
	Area = "area"

	// 生日
	Birthday = "birthday"
	// 生日时间
	BirthdayTime = "birthdayTime"

	// 性别
	Sex = "sex"
	// 顺序码
	SequenceCode = "sequenceCode"
	// 效验码
	VerifyCode = "verifyCode"
)

const (
	// 女性
	SexFemale = iota
	// 男性
	SexMale = 1
	// 18位身份证
	BitType18 = 18
	// 15位身份证
	BitType15 = 15
)
type Identity struct {
	data string
	len int
	dataMap map[string]interface{}
}


// New 新的指针指向Identity
func New(data string) *Identity {
	id := &Identity{
		data: data,
		len: len(data),
		dataMap: make(map[string]interface{}),
	}
	return id
}

// Parse 解析身份证
//
// returns 成功返回nil, 失败返回错误
func (id *Identity) Parse() (err error) {
	err = id.setBitType()
	if err != nil {
		return
	}
	err = id.setProvince()
	if err != nil {
		return
	}
	err = id.setCity()
	if err != nil {
		return
	}
	err = id.setArea()
	if err != nil {
		return
	}
	err = id.setBirthday()
	if err != nil {
		return
	}

	err = id.setSequenceCode()
	if err != nil {
		return
	}
	err = id.setVerifyCode()
	if err != nil {
		return
	}

	err = id.setSex()
	if err != nil {
		return
	}

	return
}

// setBitType 设置身份证类型
func (id *Identity) setBitType() error {
	dataLen := len(id.data)
	id.dataMap[BitType] = dataLen
	newData := ""
	switch id.dataMap[BitType] {
	case BitType18:
		newData = id.data[0:17]
		_, err := strconv.ParseUint(newData, 10, 64)
		if err != nil {
			return errors.New("identity: 18 bit identity card error")
		}
		if strings.ToUpper(string(id.data[17])) != "X" && (id.data[17]  < '0' || id.data[17] > '9') {
			return errors.New("identity: 18 bit identity card error, last bit error")
		}
	case BitType15:
		newData = id.data
		_, err := strconv.ParseUint(newData, 10, 64)
		if err != nil {
			return errors.New("identity: 15 bit identity card error")
		}
	default:
		return errors.New("identity: identity card len must be 15 or 18")
	}
	return nil
}

// GetBitType 设置身份证类型
func (id *Identity) GetBitType() interface{} {
	v := id.dataMap[BitType]
	return v
}

// setProvince 设置省份
func (id *Identity) setProvince() error {
	id.dataMap[Province] = id.data[0:2]
	return nil
}

// 获取省份
func (id *Identity) GetProvince()  interface{} {
	v := id.dataMap[Province]
	return v
}

// setCity 设置城市
func (id *Identity) setCity() error {
	id.dataMap[City] = id.data[2:4]
	return nil
}

// GetCity 获取城市
func (id *Identity) GetCity() interface{} {
	v := id.dataMap[City]
	return v
}

// setArea 设置地区
func (id *Identity) setArea() error {
	id.dataMap[Area] = id.data[4:6]
	return nil
}

// GetArea 获取地区
func (id *Identity) GetArea() interface{} {
	v := id.dataMap[Area]
	return v
}

// setBirthday 设置生日
func (id *Identity) setBirthday() error {
	switch id.GetBitType() {
	case BitType18:
		id.dataMap[Birthday] = id.data[6:14]
	case BitType15:
		id.dataMap[Birthday] = "19" + id.data[6:12]
	}
	v := id.dataMap[Birthday].(string)
	t, err := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", v[0:4], v[4:6], v[6:8]))
	if err != nil {
		return errors.New(fmt.Sprintf("identity: identity card birthday error,(birthday: %v), %#v", v, err))
	}
	id.dataMap[BirthdayTime] = t
	return nil
}

// GetBirthday 获取生日
func (id *Identity) GetBirthday() interface{} {
	v := id.dataMap[Birthday]
	return v
}

// GetBirthdayTime 获取生日
func (id *Identity) GetBirthdayTime() time.Time {
	v := id.dataMap[BirthdayTime]
	return v.(time.Time)
}

// setSex 设置性别
func (id *Identity) setSex() error {
	id.dataMap[Sex] = SexFemale
	s := id.GetSequenceCode().(string)
	if s[2] % 2 == 1 {
		id.dataMap[Sex] = SexMale
	}
	return nil
}

// GetSex 获取性别
func (id *Identity) GetSex() interface{} {
	v := id.dataMap[Sex]
	return v
}

// setSequenceCode 设置顺序码
func (id *Identity) setSequenceCode() error {
	switch id.GetBitType() {
	case BitType18:
		id.dataMap[SequenceCode] = id.data[14:17]
	case BitType15:
		id.dataMap[SequenceCode] = id.data[12:15]
	}
	return nil
}

// GetSequenceCode 获取顺序码
func (id *Identity) GetSequenceCode() interface{} {
	v := id.dataMap[SequenceCode]
	return v
}


// setVerifyCode 设置效验码
func (id *Identity) setVerifyCode() error {
	switch id.GetBitType() {
	case BitType18:
		id.dataMap[VerifyCode] = string(id.data[17])
	case BitType15:
		id.dataMap[VerifyCode] = ""
	}

	return nil
}

// GetVerifyCode 效验码
func (id *Identity) GetVerifyCode() interface{} {
	v := id.dataMap[VerifyCode]
	return v
}
