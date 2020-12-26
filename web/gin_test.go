package web

import (
	"errors"
	"github.com/hero1s/golib/task"
	"github.com/hero1s/golib/utils"
	"github.com/hero1s/golib/web/admin"
	"github.com/hero1s/golib/web/response"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hero1s/golib/log"
	"github.com/hero1s/golib/web/limit"
)

func TestInitServer(t *testing.T) {
	// 需要测试请自行解开注释测试

	c := &Config{
		Port: ":2233",
		Limit: &limit.Config{
			Rate:       0, // 0 速率不限流
			BucketSize: 100,
		},
	}

	g := InitGinServer(c)
	g.Gin.Use(LogRequest(false))
	//g.Gin.Use(JWTAuth())
	g.Register(initRoute)

	task.AddTask("test1", "0/30 * * * * *", func() error {
		log.Infof("执行定时任务:%v", time.Now().Second())
		return errors.New("test1 error")
	})
	task.AddTask("test2", "0 43 21 * * *", func() error {
		log.Infof("执行定时任务:%v", time.Now().Second())
		return nil
	})
	task.StartTask()

	admin.Init(g.Gin, "123")

	utils.RunMain(func() error {
		g.Start()
		return nil
	}, func() {

	}, "pid.txt")
}

func initRoute(g *gin.Engine) {
	cp := response.Control{}
	g.GET("/a/:abc", func(c *gin.Context) {
		log.Debugf(c.Param("abc"))
		log.Debugf(c.Request.RequestURI)
		rsp := &struct {
			Param string `json:"param"`
			Path  string `json:"path"`
		}{Param: c.Param("abc"), Path: c.Request.RequestURI}
		cp.SuccessContent(c, "成功", rsp)
	})
	g.GET("/b", func(c *gin.Context) {
		cp.Success(c, "b")
	})
	g.GET("/c", func(c *gin.Context) {
		cp.Success(c, "c")
	})
	g.GET("/d", func(c *gin.Context) {
		cp.Success(c, "d")
	})
}
