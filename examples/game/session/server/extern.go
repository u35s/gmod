package server

import "github.com/u35s/gmod"

var srv *sessionServer

func Mod() gmod.Moder {
	if srv == nil {
		srv = new(sessionServer)
	}
	return srv
}
