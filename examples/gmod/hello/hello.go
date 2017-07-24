package hello

import (
	"log"

	"github.com/u35s/gmod"
	"github.com/u35s/gmod/examples/gmod/world"
)

type hello struct {
	gmod.ModBase
	frame uint
}

func (this *hello) Run() {
	if this.frame%12 == 0 {
		this.say()
	}
	this.frame++
}

var h *hello
var Mod gmod.Moder

func init() {
	h = new(hello)
	Mod = h
}

func (this *hello) say() {
	log.Printf("hello say %v", world.Name())
}
