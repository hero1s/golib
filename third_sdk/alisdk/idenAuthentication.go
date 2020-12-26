package alisdk

import (
	"encoding/json"
	"fmt"
	"github.com/hero1s/golib/log"
	"io/ioutil"
	"net/http"
	"strings"
)

//实名认证
func IdenAuthentication(idNo, name, appCode string) (bool, error) {
	url := "http://idenauthen.market.alicloudapi.com/idenAuthentication"
	method := "POST"
	appcode := "你自己的AppCode"
	data := fmt.Sprintf("idNo=%v&name=%v", idNo, name)
	client := http.Client{}
	request, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return false, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("Authorization", "APPCODE "+appcode)
	response, err := client.Do(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}
	log.Debugf("身份验证返回:%v", string(body)) //打印返回文本
	type IdenAuthenResp struct {
		Name        string
		IdNo        string `json:"idNo"`
		RespMessage string `json:"respMessage"`
		RespCode    string `json:"respCode"`
		Province    string `json:"province"`
		City        string `json:"city"`
		County      string `json:"county"`
		Birthday    string `json:"birthday"`
		Sex         string `json:"sex"`
		Age         string `json:"age"`
	}
	var resp IdenAuthenResp
	err = json.Unmarshal(body, resp)
	if err != nil {
		return false, err
	}
	if resp.RespCode == "0000" {
		return true, nil
	}
	return false, err
}
