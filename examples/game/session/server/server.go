package server

import (
	"github.com/u35s/gmod"
	"github.com/u35s/gmod/gsrvm"
)

type sessionServer struct {
	gmod.ModBase
}

func (this *sessionServer) Init() {
	gsrvm.AddToListen(gsrvm.ToListenServer{Addr: ":8001"})
	defaultServerRoute()
	gsrvm.SetDealMsgFunc(deliverServerMsg)
}
