package agora

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"git.moumentei.com/plat_go/golib/log"
	"git.moumentei.com/plat_go/golib/third_sdk/agora/src/RtcTokenBuilder"
	"time"
)

var appID = "d3f3c1c891454bee9c7fd02449e89a0f"
var appCertificate = "86bc19d53c05442288de16e5e19bfcfc"
var appSec = "j01kmt6s"

func InitAgora(app_id, app_cert, app_sec string) {
	appID = app_id
	appCertificate = app_cert
	appSec = app_sec
}

func BuildRtcToken(uid uint64, roomid uint64) string {

	channelName := fmt.Sprintf("%v", roomid)
	expireTimeInSeconds := uint32(24 * 3600)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp := currentTimestamp + expireTimeInSeconds

	result, err := rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, uid, rtctokenbuilder.RoleAttendee, expireTimestamp)
	if err != nil {
		log.Errorf("声网创建token出错:%v", err.Error())
		return ""
	}
	log.Debugf("Token with uid: %s\n", result)
	return result
}

//检查是否是声网的抄送
func CheckChecker(AgoraSignature string, body_request string) bool {
	//hmac ,use sha1
	key := []byte(appSec)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(body_request))

	signature := hex.EncodeToString(mac.Sum(nil)[:])

	//log.Errorf("声网抄送,h:%v,b:%v",AgoraSignature,signature)
	return AgoraSignature == signature
}
