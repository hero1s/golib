package web

import (
	"fmt"
	"strings"

	"github.com/hero1s/golib/web/limit"
	"github.com/gin-gonic/gin"
)

type GinEngine struct {
	Gin  *gin.Engine
	addr string
}

func InitGinServer(c *Config) *GinEngine {
	g := gin.Default()

	if c.Limit != nil && c.Limit.Rate != 0 {
		g.Use(limit.NewLimiter(c.Limit).GinLimit())
	}

	if !strings.Contains(strings.TrimSpace(c.Port), ":") {
		c.Port = ":" + c.Port
	}

	engine := &GinEngine{Gin: g, addr: c.Host + c.Port}
	return engine
}

func (e *GinEngine) Start() {
	go func() {
		if err := e.Gin.Run(e.addr); err != nil {
			panic(fmt.Sprintf("web server addr(%s) run error(%+v).", e.addr, err))
		}
	}()
}
