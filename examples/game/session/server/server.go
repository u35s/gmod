package server

import (
	"net"

	"github.com/u35s/gmod"
	"github.com/u35s/gmod/glib/gcmd"
	"github.com/u35s/gmod/glib/gnet"
	"github.com/u35s/gmod/gsrvm"
)

type sessionServer struct {
	gmod.ModBase
}

func (this *sessionServer) Init() {
	gsrvm.AddToListen(gsrvm.ToListenServer{Addr: ":8001"})
	defaultServerRoute()
	gsrvm.SetDealMsgFunc(deliverServerMsg)
	gsrvm.SetNewAgentFunc(func(conn net.Conn) *gnet.Agent {
		return gnet.NewAgent(conn, gcmd.NewProcessor())
	})
}
