package server

import (
	"github.com/u35s/gmod"
	"github.com/u35s/gmod/gsrvm"
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
}
