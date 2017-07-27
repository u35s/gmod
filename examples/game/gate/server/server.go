package server

import (
	"net"

	"github.com/u35s/gmod"
	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/lib/gnet"
	"github.com/u35s/gmod/mods/gsrvm"
)

type gateServer struct {
	gmod.ModBase
}

func (this *gateServer) Init() {
	gsrvm.AddToConnect(gsrvm.ToConnectServer{
		Type:      "session",
		Name:      "session",
		Addr:      ":8001",
		LocalType: "gate",
		LocalName: "gate",
	})
	defaultServerRoute()
	gsrvm.SetDealMsgFunc(deliverServerMsg)
	gsrvm.SetNewAgentFunc(func(conn net.Conn) *gnet.Agent {
		return gnet.NewAgent(conn, gcmd.NewProcessor())
	})
}
