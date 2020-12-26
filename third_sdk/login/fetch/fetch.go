package fetch

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// Request request
type Request struct {
	Method string
	URL    string
	Body   []byte
	Header http.Header
}

// Cmd fetch command
func Cmd(args Request) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				deadline := time.Now().Add(30 * time.Second)
				c, err := net.DialTimeout(network, addr, 30*time.Second)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	// set request
	req, err := http.NewRequest(args.Method, args.URL,
		bytes.NewReader(
			args.Body,
		),
	)
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header = args.Header
	// get response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
