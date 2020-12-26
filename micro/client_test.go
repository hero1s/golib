package micro

import (
	"context"
	"github.com/hero1s/golib/log"
	"github.com/hero1s/golib/micro/proto"
	"github.com/micro/go-micro/v2/client"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	// 解开注释测试

	var c proto.GreeterService
	InitClient("service.client.clientName", "v1.0.0", nil, nil, func(client client.Client) {
		c = proto.NewGreeterService("service.server.serverName", client)
	})
	count := 0
	for {
		if count == 5 {
			return
		}
		time.Sleep(time.Second * 2)
		in := &proto.Request{Name: "Jerry"}
		rsp, err := c.Hello(context.Background(), in)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Debug("rsp:", rsp.Msg)
		count++
	}
}
