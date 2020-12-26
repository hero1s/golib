package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hero1s/golib/helpers/encode"
	"github.com/hero1s/golib/log"
	fetch2 "github.com/hero1s/golib/third_sdk/login/fetch"
	"strings"
)

type (
	// QQ qq
	QQ struct {
		Ret             int    `json:"ret"`            // 返回码
		Msg             string `json:"msg"`            // 如果ret<0，会有相应的错误信息提示，返回数据全部用UTF-8编码。
		NickName        string `json:"nickname"`       // 用户在QQ空间的昵称。
		FigureURL       string `json:"figureurl"`      // 大小为30×30像素的QQ空间头像URL。
		FigureURL1      string `json:"figureurl_1"`    // 大小为50×50像素的QQ空间头像URL。
		FigureURL2      string `json:"figureurl_2"`    // 大小为100×100像素的QQ空间头像URL。
		FigureURLQQ1    string `json:"figureurl_qq_1"` // 大小为40×40像素的QQ头像URL。
		FigureURLQQ2    string `json:"figureurl_qq_2"` // 大小为100×100像素的QQ头像URL。需要注意，不是所有的用户都拥有QQ的100x100的头像，但40x40像素则是一定会有。
		Gender          string `json:"gender"`         // 性别。 如果获取不到则默认返回"男"
		Sex             int64  // 值为1时是男性，值为2时是女性，值为0时是未知
		ISYellowVip     string `json:"is_yellow_vip"`      // 标识用户是否为黄钻用户（0：不是；1：是）。
		Vip             string `json:"vip"`                // 标识用户是否为黄钻用户（0：不是；1：是）
		YellowVipLevel  string `json:"yellow_vip_level"`   // 黄钻等级
		Level           string `json:"level"`              // 黄钻等级
		IsYellowYearVip string `json:"is_yellow_year_vip"` // 标识是否为年费黄钻用户（0：不是； 1：是）
	}
)
type QQToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
}

type QQAuth struct {
	QQAppID     string
	QQAppSecret string
}

func NewQQAuth(AppID, AppSecret string) *QQAuth {
	return &QQAuth{
		QQAppID:     AppID,
		QQAppSecret: AppSecret,
	}
}

//通过code来获取aceess_token及open_id
func (oAuth *QQAuth) GetQQAccessToken(code string, redirectUrl string) (string, error) {
	url := fmt.Sprintf(`https://graph.qq.com/oauth2.0/token?grant_type=authorization_code&client_id=%v&client_secret=%v&code=%v&redirect_uri=%v`,
		oAuth.QQAppID, oAuth.QQAppSecret, code, redirectUrl)
	body, err := fetch2.Cmd(fetch2.Request{
		Method: "GET",
		URL:    url,
	})
	if err != nil {
		log.Errorf("获取QQ token失败:%v   ## %v", url, err)
		return "", err
	}
	log.Debugf("获取QQ token:%v   ## %v", url, string(body))
	params := encode.ParseUrlString(string(body))
	accessToken, ok := params["access_token"]
	if !ok {
		log.Errorf("获取QQ token失败:%v  ## %v", url, string(body))
		if msg, ok := params["msg"]; ok {
			return "", errors.New(msg)
		} else {
			return "", nil
		}
	}
	return accessToken, nil
}

type QQOpenIdRep struct {
	ClientId string `json:"client_id"`
	Openid string `json:"openid"`
	Unionid string `json:"unionid"`
}

func (oAuth *QQAuth) GetQQOpenId(accessToken string) (QQOpenIdRep, error) {
	url := fmt.Sprintf(`https://graph.qq.com/oauth2.0/me?access_token=%v&unionid=1`,
		accessToken)
	body, err := fetch2.Cmd(fetch2.Request{
		Method: "GET",
		URL:    url,
	})
	var rep QQOpenIdRep
	if err != nil {
		log.Errorf("获取OpenID失败:%v", err)
		return rep, err
	}
	if strings.Contains(string(body), "callback") {
		start := strings.Index(string(body), "(")
		end := strings.LastIndex(string(body), ")")
		if start+1 > end {
			return rep, errors.New("response body error:" + string(body))
		}
		body = body[start+1 : end]
	}
	log.Debugf("获取QQ OpenID信息:%v", string(body))
	err = json.Unmarshal(body, &rep)
	if err != nil {
		log.Errorf("解析返回值错误:%v", err)
		return rep, err
	}
	return rep,nil
}

// User user
func (oAuth *QQAuth) GetQQUserInfo(accessToken, openID string) (*QQ, error) {
	var result QQ
	url := fmt.Sprintf("https://graph.qq.com/user/get_user_info?access_token=%s&oauth_consumer_key=%s&openid=%s",
		accessToken,
		oAuth.QQAppID,
		openID,
	)
	body, err := fetch2.Cmd(fetch2.Request{
		Method: "GET",
		URL:    url,
	})
	if err != nil {
		return &result, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &result, err
	}
	if result.Gender == "男" {
		result.Sex = 1
	} else {
		result.Sex = 2
	}
	if result.Ret < 0 {
		return &result, errors.New(result.Msg)
	}
	return &result, err
}
