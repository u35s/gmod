package testcmd

import "github.com/u35s/gmod/gcmd"

const CmdServer = 1

const CmdServerParam_chat = 1

type CmdServer_chat struct {
	gcmd.Cmd
	Cnt string
}

func (m *CmdServer_chat) Init() {
	m.SetBase(CmdServer,
		CmdServerParam_chat)
}
