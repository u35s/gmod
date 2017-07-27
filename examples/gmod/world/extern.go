package world

import "github.com/u35s/gmod"

func Name() string {
	return w.name()
}

var w *world

func Mod() gmod.Moder {
	if w == nil {
		w = new(world)
	}
	return w
}
