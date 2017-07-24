package server

import "github.com/u35s/gmod"

var Mod gmod.Moder
var srv *sessionServer

func init() {
	srv = new(sessionServer)
	Mod = srv
}
