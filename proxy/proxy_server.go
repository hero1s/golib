package proxy

import (
	"github.com/hero1s/golib/log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func ListenAndServe(c *Config) error {
	for i := 0; i < len(c.ProxyHost); i++ {
		up := c.ProxyHost[i]
		u, err := url.Parse(up.Upstream)
		log.Infof("{%s} => {%s}\r\n", up.Path, up.Upstream)
		if err != nil {
			log.Fatal(err.Error())
		}
		rp := httputil.NewSingleHostReverseProxy(u)
		http.HandleFunc(up.Path, func(writer http.ResponseWriter, request *http.Request) {
			if up.UpHost != "" {
				request.Host = up.UpHost
			} else {
				request.Host = u.Host
			}
			if up.TrimPath {
				request.URL.Path = strings.TrimPrefix(request.URL.Path, up.Path)
			}
			if up.IsAuth {
				auth_value := request.Header.Get("auth_key")
				if auth_value != up.AuthKey {
					writeUnAuthorized(writer)
					return
				}
			}
			rp.ServeHTTP(writer, request)
		})
	}
	// listen and serve
	log.Infof("proxy listen port:%v", c.ServerPort)
	if err := http.ListenAndServe(c.ServerPort, nil); err != nil {
		return err
	}

	return nil
}

func writeUnAuthorized(writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "Application/json")
	writer.WriteHeader(http.StatusUnauthorized)
	writer.Write([]byte("{\"status\":\"un-authorized\"}"))
}
