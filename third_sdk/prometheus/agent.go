package prometheus

import (
	"fmt"
	"git.moumentei.com/plat_go/golib/log"
	"git.moumentei.com/plat_go/golib/utils/threading"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var once sync.Once

func StartAgent(c Config) {
	once.Do(func() {
		if len(c.Host) == 0 {
			return
		}

		threading.GoSafe(func() {
			http.Handle(c.Path, promhttp.Handler())
			addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
			log.Infof("Starting prometheus agent at %s", addr)
			if err := http.ListenAndServe(addr, nil); err != nil {
				log.Error(err)
			}
		})
	})
}
