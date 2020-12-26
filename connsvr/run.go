package connsvr

import (
	"git.moumentei.com/plat_go/golib/connsvr/cluster"
	"git.moumentei.com/plat_go/golib/connsvr/module"
	"git.moumentei.com/plat_go/golib/log"
	"os"
	"os/signal"
)

//独立进程启动
func Run(mods ...module.Module) {
	log.Info("ConnSvr starting up")
	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	cluster.Init()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Infof("ConnSvr closing down (signal: %v)", sig)
	cluster.Destroy()
	module.Destroy()
}

//内部模块启动
func RunInside(end chan bool, mods ...module.Module) {
	log.Info("ConnSvr starting up")

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	cluster.Init()

	// close
	sig := <-end
	log.Infof("ConnSvr closing down (signal: %v)", sig)
	cluster.Destroy()
	module.Destroy()
}
