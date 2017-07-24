package server

import "github.com/u35s/gmod"

var Mod gmod.Moder
var srv *gateServer

func init() {
	srv = new(gateServer)
	Mod = srv
}
