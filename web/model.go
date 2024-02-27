package web

import "github.com/hero1s/golib/web/limit"

type Config struct {
	// http host
	Host string
	// http export port. :8080
	Port string
	// interface limit
	Limit      *limit.Config
	SslCrtPath string
	SslKeyPath string
}

type RecoverInfo struct {
	Time  string      `json:"time"`
	Url   string      `json:"url"`
	Err   string      `json:"error"`
	Query interface{} `json:"query"`
	Stack string      `json:"stack"`
}
