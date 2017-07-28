package gcmd

type Cmd struct {
	cmd   uint8
	param uint8
}

func (x *Cmd) SetBase(c, p uint8) { x.cmd, x.param = c, p }
func (x *Cmd) GetCmd() uint8      { return x.cmd }
func (x *Cmd) GetParam() uint8    { return x.param }

type CmdMessage struct {
	Cmd
	Data []byte
}

type Cmder interface {
	Init()
	GetCmd() uint8
	GetParam() uint8
}
