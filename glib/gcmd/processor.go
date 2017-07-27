package gcmd

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
)

const (
	maxCmdLen       = 0xffff // 最大消息长度
	compCmdLen      = 320    // 压缩消息限定
	forceCompCmdLen = 4096   // 强制压缩消息限定

	packetHeadLen      = 4
	packetHeadCompFlag = 1 << 7 //压缩标识
	packetHeadContFlag = 1 << 6 //连续标识
	packetHeadCmdMask  = 63     // 正常Cmd的Mask
)

type packetHead struct {
	Cmd   uint8
	Param uint8
	Size  uint16
}

type packetRead struct {
	packetHead
	comp    bool // 是否压缩
	cont    bool // 是否连续
	hasHead bool // 是否已读
	buf     *bytes.Buffer
}

func (x *packetRead) ParseCmd() {
	x.comp = (x.Cmd & packetHeadCompFlag) > 0
	x.cont = (x.Cmd & packetHeadContFlag) > 0
	x.Cmd = x.Cmd & packetHeadCmdMask
}

type processor struct {
	buf       bytes.Buffer
	read      packetRead
	marshal   func(interface{}) ([]byte, error)
	unmarshal func([]byte, interface{}) error
}

func (this *processor) Unmarshal(bts []byte, v interface{}) error {
	if this.unmarshal != nil {
		return this.unmarshal(bts, v)
	}
	return json.Unmarshal(bts, v)
}

func (this *processor) Marshal(i interface{}) ([]byte, error) {
	if this.marshal != nil {
		return this.marshal(i)
	}
	return json.Marshal(i)
}

func (this *processor) Pack(i interface{}) ([]byte, error) {
	if msg, ok := i.(Cmder); ok {
		msg.Init()
		bts, err := this.Marshal(msg)
		if err != nil {
			log.Printf("序列化化错误,%v", err)
			return nil, err
		}

		c := msg.GetCmd()
		if len(bts) > compCmdLen {
			/*c |= packetHeadCompFlag
			bts = glib.Compress(bts)*/
		}

		buf := new(bytes.Buffer)
		for i := 1; i > 0; {
			tmp := bts
			var cont uint8
			if len(bts) >= maxCmdLen {
				tmp = bts[:maxCmdLen]
				bts = bts[maxCmdLen:]
				cont = packetHeadContFlag
			} else {
				i = 0
			}
			binWrite(buf, c|cont)
			binWrite(buf, msg.GetParam())
			binWrite(buf, uint16(len(tmp)))
			buf.Write(tmp)
		}
		return buf.Bytes(), nil
	} else {
		return nil, errors.New("no gcmd Cmder")
	}
}

func (this *processor) UnPack(p []byte, mq chan interface{}) (n int, err error) {
	this.buf.Write(p)
	for {
		if !this.read.hasHead {
			if this.buf.Len() < packetHeadLen {
				break
			}

			binRead(&this.buf, &this.read.Cmd)
			binRead(&this.buf, &this.read.Param)
			binRead(&this.buf, &this.read.Size)
			this.read.hasHead = true
			this.read.ParseCmd()
		}
		if this.buf.Len() < int(this.read.Size) {
			break
		}
		this.read.hasHead = false

		data := this.buf.Next(int(this.read.Size))
		if this.read.cont {
			if this.read.buf == nil {
				this.read.buf = &bytes.Buffer{}
			}
			this.read.buf.Write(data)
			continue
		} else if this.read.buf != nil && this.read.buf.Len() > 0 {
			this.read.buf.Write(data)
			data = this.read.buf.Bytes()
			this.read.buf = nil
		}
		if this.read.comp {
			/*var err error
			if data, err = glib.UnCompress(data); err != nil {
				log.Printf("消息反序列化,解压缩错误:%v,%v,%v", this.read.Cmd, this.read.Param, len(data))
			}*/
		}

		msg := new(CmdMessage)
		msg.SetBase(this.read.Cmd, this.read.Param)
		msg.Data = make([]byte, len(data))
		copy(msg.Data, data)

		mq <- msg
		n++
	}
	return
}

func binRead(buf *bytes.Buffer, data interface{}) {
	binary.Read(buf, binary.LittleEndian, data)
}

func binWrite(buf *bytes.Buffer, data interface{}) {
	binary.Write(buf, binary.LittleEndian, data)
}

var marshal func(interface{}) ([]byte, error)
var unmarshal func([]byte, interface{}) error

func SetMarshal(m func(interface{}) ([]byte, error)) {
	marshal = m
}

func SetUnmarshal(u func([]byte, interface{}) error) {
	unmarshal = u
}

func NewProcessor() *processor {
	return &processor{marshal: marshal, unmarshal: unmarshal}
}
