package server

import "github.com/u35s/gmod"

var srv *gateServer

func Mod() gmod.Moder {
	if srv == nil {
		srv = new(gateServer)
	}
	return srv
}
