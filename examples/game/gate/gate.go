package main

import (
	"github.com/u35s/gmod"
	"github.com/u35s/gmod/examples/game/gate/agentm"
	"github.com/u35s/gmod/examples/game/gate/server"
	"github.com/u35s/gmod/examples/game/gate/userm"
	"github.com/u35s/gmod/mods/gsrvm"
)

func main() {
	gmod.Run(
		gsrvm.Mod(),
		server.Mod(),
		agentm.Mod(),
		userm.Mod(),
	)
}
