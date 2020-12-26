package cluster

import (
	"github.com/hero1s/golib/connsvr/chanrpc"
	"git.moumentei.com/plat_go/golib/connsvr/conf"
	"git.moumentei.com/plat_go/golib/connsvr/network"
	"git.moumentei.com/plat_go/golib/log"
	"math"
	"net"
	"time"
)

var (
	server       *network.TCPServer
	clients      []*network.TCPClient
	Processor    network.Processor
	AgentChanRPC *chanrpc.Server
)

func Init() {
	if conf.ListenAddr != "" {
		server = new(network.TCPServer)
		server.Addr = conf.ListenAddr
		server.MaxConnNum = int(math.MaxInt32)
		server.PendingWriteNum = conf.PendingWriteNum
		server.LenMsgLen = 4
		server.MaxMsgLen = math.MaxUint32
		server.NewAgent = newAgent

		server.Start()
	}

	for _, addr := range conf.ConnAddrs {
		client := new(network.TCPClient)
		client.Addr = addr
		client.ConnNum = 1
		client.ConnectInterval = 3 * time.Second
		client.PendingWriteNum = conf.PendingWriteNum
		client.LenMsgLen = 4
		client.MaxMsgLen = math.MaxUint32
		client.NewAgent = newAgent

		client.Start()
		clients = append(clients, client)
	}
}

func Destroy() {
	if server != nil {
		server.Close()
	}

	for _, client := range clients {
		client.Close()
	}
}

type agent struct {
	conn         *network.TCPConn
	processor    network.Processor
	agentChanRPC *chanrpc.Server
	userData interface{}
}

func newAgent(conn *network.TCPConn) network.Agent {
	a := new(agent)
	a.conn = conn
	return a
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debugf("read message: %v", err)
			break
		}

		if a.processor != nil {
			msg, err := a.processor.Unmarshal(data)
			if err != nil {
				log.Debugf("unmarshal message error: %v", err)
				break
			}
			err = a.processor.Route(msg, a)
			if err != nil {
				log.Debugf("route message error: %v", err)
				break
			}
		}
	}
}

func (a *agent) OnClose() {
	if a.agentChanRPC != nil {
		err := a.agentChanRPC.Call0("CloseAgent", a)
		if err != nil {
			log.Errorf("chanrpc error: %v", err)
		}
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
