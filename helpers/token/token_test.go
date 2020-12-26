package token

import (
	"git.moumentei.com/plat_go/golib/log"
	"testing"
)

func TestToken(t *testing.T) {
	tokenStr, err := GenerateToken(map[string]string{
		"uid":     "1001",
		"role_id": "122",
		"name":    "toney",
	})
	log.Info(tokenStr, err)
	custom, err := DecodeTokenByStr(tokenStr)
	log.Infof("%v,%v", custom, err)
	log.Infof("%v,%v,%v,%v,%v", custom.GetInt("uid"), custom.GetString("name"), custom.GetInt("role_id"), custom.GetString("test"), custom.GetInt("test"))

	tokenStr1, err1 := RefreshTokenByStr(tokenStr)
	log.Infof("%v,%v", tokenStr1, err1)

	custom, err = DecodeTokenByStr(tokenStr1)
	log.Infof("%v,%v", custom, err)
	log.Infof("%v,%v,%v,%v,%v", custom.GetInt("uid"), custom.GetString("name"), custom.GetInt("role_id"), custom.GetString("test"), custom.GetInt("test"))


}
