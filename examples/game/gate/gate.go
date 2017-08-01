package main

import (
	"github.com/u35s/gmod"
	"github.com/u35s/gmod/examples/game/gate/agent"
	"github.com/u35s/gmod/examples/game/gate/server"
	"github.com/u35s/gmod/examples/game/gate/user"
)

func main() {
	gmod.Run(
		server.Mod(),
		agent.Mod(),
		user.Mod(),
	)
}
