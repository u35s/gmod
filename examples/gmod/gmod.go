package main

import (
	"github.com/u35s/gmod"
	"github.com/u35s/gmod/examples/gmod/hello"
	"github.com/u35s/gmod/examples/gmod/world"
)

func main() {
	gmod.Run(
		hello.Mod,
		world.Mod,
	)
}
