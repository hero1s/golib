package web

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hero1s/golib/web/limit"
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

func (e *GinEngine) Start(c *Config) {
	go func() {
		if len(c.SslCrtPath) > 1 && len(c.SslKeyPath) > 1 {
			crtPath, _ := filepath.Abs(c.SslCrtPath)
			keyPath, _ := filepath.Abs(c.SslKeyPath)
			if err := e.Gin.RunTLS(e.addr, crtPath, keyPath); err != nil {
				panic(fmt.Sprintf("https web server addr(%s) run error(%+v).", e.addr, err))
			}
		} else {
			if err := e.Gin.Run(e.addr); err != nil {
				panic(fmt.Sprintf("http web server addr(%s) run error(%+v).", e.addr, err))
			}
		}
	}()
}
