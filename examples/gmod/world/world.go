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

func Name() string {
	return w.name()
}

var Mod gmod.Moder

var w *world

func init() {
	w = new(world)
	Mod = w
}
