package proxy

import (
	"github.com/hero1s/golib/log"
	"testing"
)

func TestService_Proxy(t *testing.T) {
	cfg, err := LoadConfig("./config.json")
	if err != nil {
		log.Fatal("reload config Err：", err)
	} else {
		if err := ListenAndServe(&cfg); err != nil {
			log.Fatal("Proxy Start Err：", err)
		}
	}
}
