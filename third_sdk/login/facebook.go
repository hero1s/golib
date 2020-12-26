package login

import (
	"encoding/json"
	"errors"
	"fmt"
	fetch2 "github.com/hero1s/golib/third_sdk/login/fetch"
)
// 这是应该web端facebook登录的，至于客户端的如何还不知道，没研究

const (
	/* 验证的过程是HTTP(s)验证，API如下：
	https://graph.facebook.com/debug_token?access_token={App-token}&input_token={User-token}
	以上的参数，User-token为用户的token，App-token为APP的token，值为 {Your AppId}%7C{Your AppSecret}。其中，%7C为urlencode的 | 符号
	*/
	facebookAuth = "https://graph.facebook.com/debug_token?input_token=%s&access_token=%s%%7C%s"
)

var (
	faceBookAppId     = ""
	facebookAppSecret = ""
)

func InitFacebookLogin(fbAppId, fbAppSecret string) {
	faceBookAppId = fbAppId
	facebookAppSecret = fbAppSecret
}
type FacebookData struct {
	Data struct {
		AppId               string   `json:"app_id"`
		Type                string   `json:"type"`
		Application         string   `json:"application"`
		DataAccessExpiresAt int64    `json:"data_access_expires_at"`
		ExpiresAt           int64    `json:"expires_at"`
		IsValid             bool     `json:"is_valid"`
		Scopes              []string `json:"scopes"`
		UserId              string   `json:"user_id"`
	} `json:"data"`
}
// facebook登录校验获取用户信息
/*
验证通过返回的Json格式（"is_valid": true),还要校验expires_at是否已经过期了
{
    "data": {
        "app_id": 000000000000000,
        "application": "Social Cafe",
        "expires_at": 1352419328,
        "is_valid": true,
        "issued_at": 1347235328,
        "scopes": [
            "email",
            "publish_actions"
        ],
        "user_id": 1207059
    }
}
*/


func VerifyFacebookToken(inputToken string) (FacebookData, error) {
	api := fmt.Sprintf(facebookAuth, inputToken, faceBookAppId, facebookAppSecret)
	data, err := fetch2.Cmd(fetch2.Request{
		Method: "GET",
		URL:    api,
		Body:   nil,
		Header: nil,
	})
	var fb FacebookData
	if err != nil {
		return fb, err
	}
	if err := json.Unmarshal(data, &fb); err != nil {
		return fb, err
	}
	if !fb.Data.IsValid {
		return fb, errors.New("facebook token invalid")
	}
	return fb, nil
}
