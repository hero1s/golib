package connsvr

import (
	"github.com/hero1s/golib/connsvr/cluster"
	"github.com/hero1s/golib/connsvr/console"
	"github.com/hero1s/golib/connsvr/module"
	"github.com/hero1s/golib/log"
	"os"
	"os/signal"
	"sync"
)

var (
	endChan chan bool
	mu      sync.Mutex
	once    sync.Once
)

// 独立进程启动
func Run(mods ...module.Module) {
	log.Info("ConnSvr starting up")
	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	cluster.Init()

	// console
	console.Init()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	defer func() {
		sig := <-c
		log.Infof("ConnSvr closing down (signal: %v)", sig)

		console.Destroy()
		cluster.Destroy()
		module.Destroy()
	}()

	select {}
}

// 内部模块启动
func RunInside(mods ...module.Module) {
	log.Info("ConnSvr starting up")

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	cluster.Init()

	// console
	console.Init()
	once.Do(func() {
		endChan = make(chan bool, 1)
	})
	// close
	sig := <-endChan
	log.Infof("ConnSvr closing down (signal: %v)", sig)
	console.Destroy()
	cluster.Destroy()
	module.Destroy()
}
func Stop() {
	mu.Lock()
	defer mu.Unlock()
	if endChan != nil {
		endChan <- true
	}
}
