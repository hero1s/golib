package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hero1s/golib/log"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

var accessKeyId = ""
var accessSecret = ""
var regionId = "cn-hangzhou"

func InitAliPhoneLogin(keyId, secret string) {
	accessKeyId = keyId
	accessSecret = secret
}

func AliGetMobile(accessToken string) (string, error) {
	client, err := sdk.NewClientWithAccessKey(regionId, accessKeyId, accessSecret)
	if err != nil {
		log.Errorf("阿里登录初始化错误:%v", err)
		return "", err
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dypnsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "GetMobile"
	request.QueryParams["RegionId"] = regionId
	request.QueryParams["AccessToken"] = accessToken

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.Errorf("阿里登录请求错误:%v", err)
		return "", err
	}
	repStr := response.GetHttpContentString()
	log.Debugf("阿里登录返回:", repStr)
	type GetMobileResultDTO struct {
		Mobile string `json:"Mobile"`
	}
	type AliRep struct {
		RequestId string             `json:"RequestId"`
		Code      string             `json:"Code"`
		Message   string             `json:"Message"`
		PhoneInfo GetMobileResultDTO `json:"GetMobileResultDTO"`
	}
	var rep AliRep
	err = json.Unmarshal([]byte(repStr), &rep)
	if err != nil {
		log.Errorf("解析返回失败:%v", err)
		return "", err
	}
	if rep.Code != "OK" {
		log.Errorf("阿里手机登录错误:%+v", rep)
		return "", errors.New(fmt.Sprintf("阿里手机登录异常:%v", rep.Message))
	}
	log.Debugf("返回结果:%#v", rep)
	return rep.PhoneInfo.Mobile, nil
}
