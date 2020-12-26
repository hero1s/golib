package login

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	fetch2 "github.com/hero1s/golib/third_sdk/login/fetch"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// 第三方登录校验
const (
	// google的直接拿token去获取用户的信息，验证
	googleAuth = "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s"
)

// 相关appId,appSecret
var (
	googleAppId     = ""
	googleAppSecret = ""
)

func InitGoogleLogin(goAppId, goAppSecret string) {
	googleAppId = goAppId
	googleAppSecret = goAppSecret
}


type GoogleData struct {
	Iss           string `json:"iss"`
	Sub           string `json:"sub"`
	Azp           string `json:"azp"`
	Email         string `json:"email"`
	AtHash        string `json:"at_hash"`
	EmailVerified bool   `json:"email_verified"`
	Aud           string `json:"aud"`
	Iat           string `json:"iat"`
	Exp           string `json:"exp"`
	Name          string `json:"name"`
	Picture       string `json:"picture"` //用户头像
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
}

// google登录校验获取用户信息
/*
验证通过返回的Json格式（校验aud是否与后台保存的AppId匹配）还要校验exp是否已经过期了：
{
 "iss": "https://accounts.google.com",
 "sub": "110169484474386276334",  //用户id
 "azp": "1008719970978-hb24n2dstb40o45d4feuo2ukqmcc6381.apps.googleusercontent.com",
 "email": "billd1600@gmail.com",
 "at_hash": "X_B3Z3Fi4udZ2mf75RWo3w",
 "email_verified": "true",
 "aud": "1008719970978-hb24n2dstb40o45d4feuo2ukqmcc6381.apps.googleusercontent.com",  // appId
 "iat": "1433978353",
 "exp": "1433981953"
 "name" : "Test User",
 "picture": "https://lh4.googleusercontent.com/-kYgzyAWpZzJ/ABCDEFGHI/AAAJKLMNOP/tIXL9Ir44LE/s99-c/photo.jpg",
 "given_name": "Test",
 "family_name": "User",
 "locale": "en"
}
*/
func GoogleLogin(idToken string) (bool, GoogleData, error) {
	api := fmt.Sprintf(googleAuth, idToken)
	var goData GoogleData
	data, err := fetch2.Cmd(fetch2.Request{
		Method: "GET",
		URL:    api,
		Body:   nil,
		Header: nil,
	})
	if err != nil {
		return false, goData, err
	}
	if err := json.Unmarshal(data, &goData); err != nil {
		return false, goData, err
	}
	if goData.Iss != "accounts.google.com" || goData.Iss != "https://accounts.google.com" {
		return false, goData, nil
	}
	if goData.Aud != googleAppId {
		return false, goData, nil
	}
	return true, goData, nil
}

func VerifyGoogleIdToken(idToken string) (*oauth2.Tokeninfo, error) {
	ctx := context.Background()
	oauth2Service, err := oauth2.NewService(ctx, option.WithoutAuthentication())
	if err != nil {
		return nil, err
	}
	tokenInfo, err := oauth2Service.Tokeninfo().IdToken(idToken).Do()
	if err != nil {
		return tokenInfo, err
	}
	if tokenInfo.Audience != googleAppId {
		return tokenInfo, errors.New("app_id不匹配,不是合法的id_token")
	}
	return tokenInfo, nil
}
