package web

import (
	"errors"
	"github.com/hero1s/golib/task"
	"github.com/hero1s/golib/utils"
	"github.com/hero1s/golib/web/admin"
	"github.com/hero1s/golib/web/middleware"
	"github.com/hero1s/golib/web/response"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hero1s/golib/log"
	lconf "github.com/hero1s/golib/log/conf"
	"github.com/hero1s/golib/web/limit"
)

func TestInitServer(t *testing.T) {
	// 需要测试请自行解开注释测试
	initLog(true)
	log.Infof("start test init server")
	log.Errorf("test lark msg")
	c := &Config{
		Port: ":2233",
		Limit: &limit.Config{
			Rate:       0, // 0 速率不限流
			BucketSize: 100,
		},
	}

	g := InitGinServer(c)
	g.Gin.Use(middleware.LogRequest(nil, nil))
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
		g.Start(c)
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

func initLog(isDev bool) {
	log.InitLogger(
		lconf.WithProjectName(""),
		lconf.WithIsStdOut(!isDev),
		lconf.WithFilename("logs/api_svr/api.log"),
		lconf.WithMaxAge(10),
		lconf.WithMaxSize(50), lconf.WithProjectName("test"))
	log.SetLarkUrl("https://open.feishu.cn/open-apis/bot/v2/hook/09229461-36e5-4cd6-be8d-e55d3eb794a8")
}
