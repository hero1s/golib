package gate

import (
	"github.com/hero1s/golib/connsvr/chanrpc"
	"github.com/hero1s/golib/connsvr/gate"
	"github.com/hero1s/golib/connsvr/network"
	"time"
)

//  实现一个简易的gate 模块

type Module struct {
	*gate.Gate
}

var GateModule *Module

func InitModule(WSAddr string, TCPAddr string, agentChanRPC *chanrpc.Server, processor network.Processor) *Module {
	GateModule = new(Module)
	GateModule.Init(WSAddr, TCPAddr, agentChanRPC, processor)
	return GateModule
}

func (m *Module) Init(WSAddr string, TCPAddr string, agentChanRPC *chanrpc.Server, processor network.Processor) {
	m.Gate = &gate.Gate{
		MaxConnNum:      10000,
		PendingWriteNum: 2000,
		MaxMsgLen:       8192,
		WSAddr:          WSAddr,
		HTTPTimeout:     10 * time.Second,
		TCPAddr:         TCPAddr,
		LenMsgLen:       2,
		LittleEndian:    false,
		AgentChanRPC:    agentChanRPC,
		Processor:       processor,
	}
}

func (m *Module) OnInit() {

}
