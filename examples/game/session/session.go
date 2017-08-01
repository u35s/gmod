package main

import (
	"github.com/u35s/gmod"
	"github.com/u35s/gmod/examples/game/session/server"
)

func main() {
	gmod.Run(
		server.Mod(),
	)
}
