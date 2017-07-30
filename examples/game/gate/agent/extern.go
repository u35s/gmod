package agent

import "github.com/u35s/gmod"

var agentm *agentManager

func Mod() gmod.Moder {
	if agentm == nil {
		agentm = new(agentManager)
	}
	return agentm
}
