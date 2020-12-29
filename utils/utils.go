package utils

import (
	"github.com/hero1s/golib/helpers/file"
	"github.com/hero1s/golib/log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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
	//截获dump信息输出
	file, err1 := file.Open(pidFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err1 != nil {
		log.Errorf("pid文件打开失败")
	} else {
		syscall.Dup2(int(file.Fd()), 1)
		syscall.Dup2(int(file.Fd()), 2)
	}

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Infof("server closing down (signal: %v)", sig)
	destroy()
	return nil
}
