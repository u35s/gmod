package testcmd

import "github.com/u35s/gmod/lib/gcmd"

const CmdUser = 2

const CmdUserParam_login = 1

type CmdUser_login struct {
	gcmd.Cmd
	Seqid uint
	Accid uint
}

func (m *CmdUser_login) Init() {
	m.SetBase(CmdUser,
		CmdUserParam_login)
}

const CmdUserParam_chat = 2

type CmdUser_chat struct {
	gcmd.Cmd
	Cnt string
}

func (m *CmdUser_chat) Init() {
	m.SetBase(CmdUser,
		CmdUserParam_chat)
}
