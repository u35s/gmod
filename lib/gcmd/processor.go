package gcmd

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	packetHeadLen = 2
)

type processor struct {
	buf bytes.Buffer
}

func (this *processor) Unmarshal(bts []byte) (interface{}, error) {
	return this.unmarshal(bts)
}

func (this *processor) unmarshal(bts []byte) (*CmdMessage, error) {
	if len(bts) < packetHeadLen {
		return nil, errors.New("bts too small")
	}
	this.buf.Write(bts)

	msg := new(CmdMessage)
	binRead(&this.buf, &msg.cmd)
	binRead(&this.buf, &msg.param)
	msg.Data = this.buf.Bytes()
	this.buf.Reset()
	return msg, nil
}

func (this *processor) Marshal(m interface{}) ([]byte, error) {
	if msg, ok := m.(Cmder); ok {
		return this.marshal(msg)
	} else {
		return nil, errors.New(fmt.Sprintf("%v does not implement the gcmd.Cmder interface", m))
	}
}

func (this *processor) marshal(msg Cmder) ([]byte, error) {
	msg.Init()
	bts, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	binWrite(buf, msg.GetCmd())
	binWrite(buf, msg.GetParam())
	binWrite(buf, bts)
	return buf.Bytes(), nil
}

func binRead(buf *bytes.Buffer, data interface{}) {
	binary.Read(buf, binary.LittleEndian, data)
}

func binWrite(buf *bytes.Buffer, data interface{}) {
	binary.Write(buf, binary.LittleEndian, data)
}

func NewProcessor() *processor {
	return &processor{}
}
