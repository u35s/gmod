package testcmd

import "github.com/u35s/gmod/lib/gcmd"

const CmdServer = 1

const CmdServerParam_establishConnection = 1

type CmdServer_establishConnection struct {
	gcmd.Cmd
	Type string
	Name string
}

func (m *CmdServer_establishConnection) Init() {
	m.SetBase(CmdServer,
		CmdServerParam_establishConnection)
}

const CmdServerParam_forwardUserMsg = 2

type CmdServer_forwardUserMsg struct {
	gcmd.Cmd
	UID      uint
	SubCmd   uint8
	SubParam uint8
	Data     []byte
}

func (m *CmdServer_forwardUserMsg) Init() {
	m.SetBase(CmdServer,
		CmdServerParam_forwardUserMsg)
}

const CmdServerParam_userLogin = 3

type CmdServer_userLogin struct {
	gcmd.Cmd
	Seqid uint
	Accid uint
}

func (m *CmdServer_userLogin) Init() {
	m.SetBase(CmdServer,
		CmdServerParam_userLogin)
}

const CmdServerParam_chat = 4

type CmdServer_chat struct {
	gcmd.Cmd
	Cnt string
}

func (m *CmdServer_chat) Init() {
	m.SetBase(CmdServer,
		CmdServerParam_chat)
}
