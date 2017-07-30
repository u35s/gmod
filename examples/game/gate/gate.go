package main

import (
	"github.com/u35s/gmod"
	"github.com/u35s/gmod/examples/game/gate/agent"
	"github.com/u35s/gmod/examples/game/gate/server"
	"github.com/u35s/gmod/examples/game/gate/user"
	"github.com/u35s/gmod/mods/gsrvs"
)

func main() {
	gmod.Run(
		gsrvs.Mod(),
		server.Mod(),
		agent.Mod(),
		user.Mod(),
	)
}
