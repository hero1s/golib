package http_client

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Method string
	URL    string
	Body   string
	Header http.Header
}

func HttpsRequest(args Request) ([]byte, error) {
	client := &http.Client{Transport: &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			deadline := time.Now().Add(20 * time.Second)
			c, err := net.DialTimeout(network, addr, 18*time.Second)
			if err != nil {
				return nil, err
			}
			c.SetDeadline(deadline)
			return c, nil
		},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
	}
	req, err := http.NewRequest(args.Method, args.URL, strings.NewReader(args.Body))
	if err != nil {
		return nil, nil
	}
	req.Close = true
	req.Header = args.Header

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//发送钉钉消息
func SendMsgToDingDing(msg string, tokenUrl string, at []string) {
	needAtAll := false
	if at == nil {
		needAtAll = true
	}
	headers := http.Header{}
	headers.Add("Content-Type", "application/json;charset=utf-8")
	data := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": msg,
		},
		"at": map[string]interface{}{
			"atMobiles": at,
			"isAtAll":   needAtAll,
		},
	}
	buff, _ := json.Marshal(data)
	HttpsRequest(Request{
		Method: "POST",
		URL:    tokenUrl,
		Header: headers,
		Body:   string(buff),
	})
}
