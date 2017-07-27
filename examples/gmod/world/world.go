package world

import (
	"github.com/u35s/gmod"
)

type world struct {
	gmod.ModBase
}

func (this *world) name() string {
	return "world"
}
