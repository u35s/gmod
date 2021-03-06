package gcmd

import (
	"shell/cmd"
	"testing"
)

type Cmd_test struct {
	cmd.Cmd
	Data uint
}

func (m *Cmd_test) Init() {
	m.SetBase(1, 2)
}

func Test_PackDrawCmd(t *testing.T) {
	processor := NewProcessor()

	var msg Cmd_test
	msg.Init()
	msg.Data = 1
	bts, _ := processor.Marshal(&msg)
	t.Log(string(bts))
	m, _ := processor.Unmarshal(bts)
	log := t.Log
	if m.GetCmd() != 1 || m.GetParam() != 2 {
		log = t.Error
	}
	log(m.GetCmd(), m.GetParam(), string(m.Data))

}
