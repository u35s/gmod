package main

import (
	"github.com/u35s/gmod"
	"github.com/u35s/gmod/examples/game/session/server"
	"github.com/u35s/gmod/mods/gsrvm"
)

func main() {
	gmod.Run(
		gsrvm.Mod(),
		server.Mod(),
	)
}
