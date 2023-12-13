package web

import (
	"github.com/gin-gonic/gin"
	"github.com/hero1s/golib/web/middleware"
)

type HandlerRegisterFunc func(router *gin.Engine)

// 路由总入口，注册所有API的 路由
func (g *GinEngine) Register(registerFunc HandlerRegisterFunc) {
	//配置跨域
	g.Gin.Use(middleware.CORS(), middleware.Recovery(true), middleware.RequestId())
	//g.Gin.HandleMethodNotAllowed = true

	// 注册router
	if registerFunc != nil {
		registerFunc(g.Gin)
	}
}
