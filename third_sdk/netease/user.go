package netease

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hero1s/golib/log"
)

//--------------------网易云信ID----------------------------------------

func (b *Net) CreateUserId(data map[string]interface{}) (r CreateUserID, err error) {
	url := "https://api.netease.im/nimserver/user/create.action"
	rsp, err := b.postDataHttps(url, data)
	log.Debugf("创建云信id,post返回:%v",string(rsp[:]))
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return r, err
	}

	if r.Code != 200 {
		return r, errors.New(fmt.Sprintf("%v", r.Code))
	}
	return r, nil
}

func (b *Net) UpdateUserId(data map[string]interface{}) error {
	url := "https://api.netease.im/nimserver/user/update.action"
	return handleOnlyCodeResponse(b.postDataHttps(url, data))
}
func (b *Net) UpdateUserToken(accId string) (string, error) {
	url := "https://api.netease.im/nimserver/user/refreshToken.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{"accid": accId})
	log.Debugf("更新云信用户token,post返回%v",string(rsp[:]))
	if err != nil {
		return "", err
	}
	var r RefreshToken
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return "", err
	}
	if r.Code != 200 {
		log.Debugf("更新token,返回:%#v",r)
		return "",err
	}
	return r.Info.Token, nil
}

func (b *Net) BlockUserId(accId string, needKick bool) error {
	url := "https://api.netease.im/nimserver/user/block.action"
	return handleOnlyCodeResponse(b.postDataHttps(url, map[string]interface{}{"accid": accId, "needkick": needKick}))
}

func (b *Net) UnblockUserId(accId string) error {
	url := "https://api.netease.im/nimserver/user/unblock.action"
	return handleOnlyCodeResponse(b.postDataHttps(url, map[string]interface{}{"accid": accId}))
}

//--------------------用户名片----------------------------------------

func (b *Net) UpdateUserInfo(data map[string]interface{}) error {
	url := "https://api.netease.im/nimserver/user/updateUinfo.action"
	return handleOnlyCodeResponse(b.postDataHttps(url, data))
}

func (b *Net) GetUserInfo(accids []string) (r GetUserInfo, err error) {
	url := "https://api.netease.im/nimserver/user/getUinfos.action"
	accidStr, _ := json.Marshal(accids)
	rsp, err := b.postDataHttps(url, map[string]interface{}{"accids": string(accidStr)})
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return r, err
	}
	if r.Code != 200 {
		return r, errors.New(fmt.Sprintf("%v", r.Code))
	}
	return r, nil

}

//--------------------用户设置----------------------------------------

func (b *Net) SetNotifyDevice(accId, donnopOpen bool) error {
	url := "https://api.netease.im/nimserver/user/setDonnop.action"
	return handleOnlyCodeResponse(b.postDataHttps(url, map[string]interface{}{"accid": accId, "donnoOpen": donnopOpen}))
}
