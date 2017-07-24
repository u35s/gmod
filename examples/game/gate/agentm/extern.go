package agentm

import "github.com/u35s/gmod"

var agentm *agentManager

var Mod gmod.Moder

func init() {
	agentm = new(agentManager)
	Mod = agentm
}
