package netease

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"github.com/hero1s/golib/log"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	neturl "net/url"
	"strings"
	"time"
)

//mnetease IM SDK
type Net struct {
	AppKey    string
	AppSecret string
	Nonce     string
	CurTime   string
	CheckSum  string
}

var (
	AppKey    = "ee931139c93a6750dd6e467c9e82f5a4"
	AppSecret = "774f1b1ec132"
)

func NewNet(appKey, appSecret string) *Net {
	return &Net{
		AppKey:    appKey,
		AppSecret: appSecret,
	}
}

func (b *Net) checkSumBuilder() {
	charHex := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	nonce := ""
	for i := 0; i < 64; i++ {
		nonce = nonce + charHex[rand.Intn(16)]
	}
	b.Nonce = nonce
	b.CurTime = fmt.Sprintf("%v", time.Now().Unix())
	CheckStr := b.AppSecret + b.Nonce + b.CurTime
	h := sha1.New()
	h.Write([]byte(CheckStr))
	b.CheckSum = hex.EncodeToString(h.Sum(nil))
}

//检查是否是云信的抄送
func (b *Net) CheckSumChecker(checkSum_request string, CurTime_request string, body_request string, md5_requist string) bool {
	m := md5.New()
	m.Write([]byte(body_request))
	md5_local := hex.EncodeToString(m.Sum(nil))
	if md5_local != md5_requist {
		log.Errorf("检查是否是云信的抄送,md5错误,local:%v,request:%v", md5_local, md5_requist)
		return false
	}

	sha_str := b.AppSecret + md5_local + CurTime_request
	h := sha1.New()
	h.Write([]byte(sha_str))
	checkSum_local := hex.EncodeToString(h.Sum(nil))

	if checkSum_local != checkSum_request {
		log.Errorf("检查是否是云信的抄送,sha错误,local:%v,request:%v", checkSum_local, checkSum_request)
		return false
	}
	return true
}

func (b *Net) getDataHttps(url string, data map[string]interface{}) ([]byte, error) {
	b.checkSumBuilder()
	headers := http.Header{
		"AppKey":       {b.AppKey},
		"Nonce":        {b.Nonce},
		"CurTime":      {b.CurTime},
		"CheckSum":     {b.CheckSum},
		"Content-Type": {"application/x-www-form-urlencoded;charset=utf-8"},
	}
	var dataStr string
	for k, v := range data {
		if dataStr == "" {
			dataStr = fmt.Sprintf("%v=%v", k, v)
		} else {
			dataStr = dataStr + fmt.Sprintf("&%v=%v", k, v)
		}
	}
	url = fmt.Sprintf("%v%v", url, dataStr)
	log.Debugf("云信get请求:%v", url)
	rsp, err := httpsRequest(Request{
		Method: "GET",
		URL:    url,
		Header: headers,
		//Body:   dataStr,
	})
	log.Debugf("云信get请求返回:%v", string(rsp[:]))
	return rsp, err
}

func (b *Net) postDataHttps(url string, data map[string]interface{}) ([]byte, error) {
	b.checkSumBuilder()
	headers := http.Header{
		"AppKey":       {b.AppKey},
		"Nonce":        {b.Nonce},
		"CurTime":      {b.CurTime},
		"CheckSum":     {b.CheckSum},
		"Content-Type": {"application/x-www-form-urlencoded;charset=utf-8"},
	}
	tmp := neturl.Values{}
	for k, v := range data {
		tmp.Add(k, fmt.Sprintf("%v", v))
	}
	dataStr := tmp.Encode()
	rsp, err := httpsRequest(Request{
		Method: "POST",
		URL:    url,
		Header: headers,
		Body:   dataStr,
	})
	log.Debugf("云信post请求返回:%v", string(rsp[:]))
	return rsp, err
}

type Request struct {
	Method string
	URL    string
	Body   string
	Header http.Header
}

func httpsRequest(args Request) ([]byte, error) {
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
	log.Debugf("body:%v", args.Body)
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
