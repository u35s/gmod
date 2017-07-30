package main

import (
	"github.com/u35s/gmod"
	"github.com/u35s/gmod/examples/game/session/server"
	"github.com/u35s/gmod/mods/gsrvs"
)

func main() {
	gmod.Run(
		gsrvs.Mod(),
		server.Mod(),
	)
}
