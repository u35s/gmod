package hello

import "github.com/u35s/gmod"

var h *hello

func Mod() gmod.Moder {
	if h == nil {
		h = new(hello)
	}
	return h
}
