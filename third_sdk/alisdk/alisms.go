package alisdk

import (
	"errors"
	"fmt"
	"github.com/hero1s/golib/cache"
	"github.com/hero1s/golib/helpers/math"
	"github.com/denverdino/aliyungo/sms"
	"strconv"
	"time"
)

var (
	gatewayUrl         = "http://dysmsapi.aliyuncs.com/"
	accessKeyIdSMS     = "LTAIPT1E4BInyCX5"
	accessKeySecretSMS = "OXsgdCInqzbhxbDMVhwMJylzojyhTm"
)

func InitAliSMS(accessKey, accessKeySecret string) {
	accessKeyIdSMS = accessKey
	accessKeySecretSMS = accessKeySecret
}

//发送身份验证短信
func SendSMSCaptcha(phone, signName, templateCode, phoneKey string, c cache.Cache, timeout, sendAgainTimeout int64) (string, string, error) {
	sendAgainKey := fmt.Sprintf("%v%v-%v", phoneKey, phone, "SendAgain")
	if c.IsExist(sendAgainKey) {
		return "", "短信验证码已发送,无需重复发送", nil
	}
	code := int64(math.Random(100000, 999999))
	templateParam := fmt.Sprintf("{\"code\":\"%d\"}", code)
	client := sms.NewDYSmsClient(accessKeyIdSMS, accessKeySecretSMS)
	args := sms.SendSmsArgs{
		SignName:      signName,
		TemplateCode:  templateCode,
		TemplateParam: templateParam,
		PhoneNumbers:  phone,
	}
	resp, err := client.SendSms(&args)
	if err != nil {
		return "", resp.Message, err
	}
	captcha := strconv.FormatInt(code, 10)
	err = c.Put(phoneKey+phone, captcha, time.Duration(timeout)*time.Second)
	err = c.Put(sendAgainKey, captcha, time.Duration(sendAgainTimeout)*time.Second)
	return captcha, resp.Message, err
}

//检查验证码
func VerifySMSCaptcha(phone, captcha, phoneKey string, c cache.Cache) (error, bool) {
	sendAgainKey := fmt.Sprintf("%v%v-%v", phoneKey, phone, "SendAgain")
	cacheCaptcha := c.GetString(phoneKey + phone)
	if cacheCaptcha == "" {
		return errors.New("验证码已过期"), false
	}
	if cacheCaptcha != captcha {
		return errors.New(fmt.Sprintf("验证码校验不通过:%s--%s",cacheCaptcha,captcha)), false
	}
	c.Delete(sendAgainKey) //验证通过删除再次发送时间限制
	return nil, true
}

func IsExistSMSCaptcha(phoneNum, phoneKey string, c cache.Cache) bool {
	return c.IsExist(phoneKey + phoneNum)
}

//发送短信通知
func SendSMSNotify(phone, signName, templateCode, phoneKey,notifyMsg string, c cache.Cache, timeout, sendAgainTimeout int64) (string, error) {
	sendAgainKey := fmt.Sprintf("%v%v-%v", phoneKey, phone, "SendAgain")
	if c.IsExist(sendAgainKey) {
		return "短信验证码已发送,无需重复发送", nil
	}
	templateParam := fmt.Sprintf("{\"name\":\"%v\"}", notifyMsg)
	client := sms.NewDYSmsClient(accessKeyIdSMS, accessKeySecretSMS)
	args := sms.SendSmsArgs{
		SignName:      signName,
		TemplateCode:  templateCode,
		TemplateParam: templateParam,
		PhoneNumbers:  phone,
	}
	resp, err := client.SendSms(&args)
	if err != nil {
		return resp.Message, err
	}
	err = c.Put(phoneKey+phone, notifyMsg, time.Duration(timeout)*time.Second)
	err = c.Put(sendAgainKey, notifyMsg, time.Duration(sendAgainTimeout)*time.Second)
	return resp.Message, err
}