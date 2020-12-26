package utils

import (
	"git.moumentei.com/plat_go/golib/helpers/file"
	"git.moumentei.com/plat_go/golib/log"
	"os"
	"os/signal"
	"runtime"
)

// 获取正在运行的函数名
func RunFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
func CallerFuncName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}

// 标准启动
func RunMain(initFunc func() error, destroy func(), pidFileName string) error {
	log.Info("server start")
	if err := initFunc(); err != nil {
		log.Errorf("server start error:%v", err.Error())
		return err
	}
	if len(pidFileName) > 1 {
		file.WritePidFile(pidFileName)
	}
	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Infof("server closing down (signal: %v)", sig)
	destroy()
	return nil
}
