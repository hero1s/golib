package webhook

import (
	"encoding/json"
	"github.com/hero1s/golib/log"
	"github.com/hero1s/golib/web/xhttp"
	"net/http"
)

// 飞书
func SendLark(msg string, tokenUrl string) (bool, error) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json;charset=utf-8")
	data := map[string]interface{}{
		"msg_type": "text",
		"content":  map[string]string{"text": msg},
	}
	buff, _ := json.Marshal(data)
	resp, err := xhttp.Request("POST", tokenUrl, xhttp.WithBodyString(string(buff)), xhttp.WithHeader(headers))
	if err != nil {
		log.Errorf("发送飞书消息失败:%v,%v", resp.StatusCode, err)
		return false, err
	}
	return true, nil
}

// 钉钉
func SendMsgToDingDing(msg string, tokenUrl string) (bool, error) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json;charset=utf-8")
	type Msg struct {
		Content string `json:"content"`
	}
	var m Msg
	m.Content = msg
	data := map[string]interface{}{
		"msgtype": "text",
		"text":    m,
	}
	buff, _ := json.Marshal(data)
	resp, err := xhttp.Request("POST", tokenUrl, xhttp.WithBodyString(string(buff)), xhttp.WithHeader(headers))
	if err != nil {
		log.Errorf("发送钉钉消息失败:%v,%v", resp.StatusCode, err)
		return false, err
	}
	return true, nil
}

// 企业微信
func SendWechat(msg string, tokenUrl string) (bool, error) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json;charset=utf-8")
	type Msg struct {
		Content string `json:"content"`
	}
	var m Msg
	m.Content = msg
	data := map[string]interface{}{
		"msgtype": "text",
		"text":    m,
	}
	buff, _ := json.Marshal(data)
	resp, err := xhttp.Request("POST", tokenUrl, xhttp.WithBodyString(string(buff)), xhttp.WithHeader(headers))
	if err != nil {
		log.Errorf("发送钉钉消息失败:%v,%v", resp.StatusCode, err)
		return false, err
	}
	return true, nil
}
