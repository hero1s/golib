package netease

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNet_CheckSum(t *testing.T) {
	net := NewNet(AppKey, AppSecret)
	net.checkSumBuilder()
	ex := struct {
		Constellation string
	}{
		Constellation: "水平座",
	}
	exJson, _ := json.Marshal(ex)
	data := map[string]interface{}{
		"accid": "lklk6nj95wtio09_112345",
		"name":  "娟",
		//"props":  `{"sex":1,"man":2}`,
		"icon": "http://hugapp.oss-cn-shenzhen.aliyuncs.com/platform/internal_test/101441/file/1530291698876911638101441.png",
		"sign": "",
		//"email":  "dghpgyss@163.com",
		"birth":  "2011-01-26",
		"mobile": "17374414529",
		"gender": 2,
		"ex":     string(exJson),
	}
	rsp, err := net.CreateUserId(data)
	if err != nil {
		fmt.Println("-----网易云信返回的错误码::", err)
	}
	fmt.Printf("-----网易云信的返回:%#v\n", rsp)
	data1 := map[string]interface{}{
		"accid":  "12345678",
		"name":   "test name",
		"props":  "{sex:1,man:2}",
		"icon":   "imageurl",
		"sign":   "test sign",
		"email":  "dghpgyss@163.com",
		"birth":  "2018-05-26",
		"mobile": "13712345678",
		"gender": 1,
		"ex":     "just ex",
	}
	rsp, err = net.CreateUserId(data1)
	if err != nil {
		fmt.Println("-----err:", err)
	}
	fmt.Println("-----rsp 2:", rsp)

	/*
		data2 := map[string]interface{}{
			"accid":  "123456789",
			"name":   "test name",
			"icon":   "imageurl---------asdfk",
			"sign":   "test sign",
			"email":  "dghpgyss@163.com",
			"birth":  "2018-05-26",
			"mobile": "13712345678",
			"gender": 1,
		}
		err = net.UpdateUserInfo(data2)
		if err != nil {
			fmt.Println("-----err:", err)
		}
		fmt.Println("-----rsp 3:", rsp)
	*/

}

/*
func TestNet_UpdateUserToken(t *testing.T) {
	net := NewNet(AppKey, AppSecret)
	net.checkSumBuilder()
	rsp, err := net.UpdateUserToken("119iiiii171")
	if err != nil {
		fmt.Println("-----err:", err)
		return
	}
	fmt.Println("-----rsp:", rsp)
}
fweasblywlsmmtc_102143
*/

func TestNet_GetUserInfo(t *testing.T) {
	fmt.Println("---------------------下面是获取用户信息-------------------------")
	net := NewNet(AppKey, AppSecret)
	net.checkSumBuilder()
	//accids := []string{"123456", "12345678", "123456789"}
	accids := []string{"rwgwvevjffswmyq_147223", "123456", "lklk6nj95wtio09_112345"}
	//accids := []string{"abcdefgddaa123456", "123456"}
	rsp, err := net.GetUserInfo(accids)
	if err != nil {
		fmt.Println("-----err:", err)
		return
	}
	fmt.Printf("GetUserInfo-rsp:%#v\n", rsp)
}

func TestNet_GetFriend(t *testing.T) {
	fmt.Println("---------------------下面是获取用户好友-------------------------")
	net := NewNet(AppKey, AppSecret)
	net.checkSumBuilder()
	rsp, err := net.GetFriend("fweasblywlsmmtc_102143", 0)
	if err != nil {
		fmt.Println("-----err:", err)
		return
	}
	fmt.Printf("GetFriend-rsp:%#v\n", rsp)
}
