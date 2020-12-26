package sdk

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	timeout = time.Second * 10
	client  = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConnsPerHost: 1024,
		},
		Timeout: timeout,
	}
)

// Get 发起Get请求
func Get(uri string, params, header map[string]string) ([]byte, error) {
	return doHTTP(http.MethodGet, uri, params, header)
}

// Post 发起Post请求
func Post(uri string, params, header map[string]string) ([]byte, error) {
	return doHTTP(http.MethodPost, uri, params, header)
}

// doHTTP 发起http请求
func doHTTP(method, uri string, params, header map[string]string) (data []byte, err error) {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	reader := values.Encode()
	if method == http.MethodGet {
		uri = fmt.Sprintf("%s?%s", uri, values.Encode())
		reader = ""
	}

	req, err := http.NewRequest(method, uri, strings.NewReader(reader))
	if err != nil {
		return
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("wrong response status:%d", resp.StatusCode)
		return
	}

	data, err = ioutil.ReadAll(resp.Body)
	return
}
