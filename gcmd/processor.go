package gcmd

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
)

const (
	MaxCmdLen       = 0xffff // 最大消息长度
	CompCmdLen      = 320    // 压缩消息限定
	ForceCompCmdLen = 4096   // 强制压缩消息限定

	PacketHeadLen      = 4
	PacketHeadCompFlag = 1 << 7 //压缩标识
	PacketHeadContFlag = 1 << 6 //连续标识
	PacketHeadCmdMask  = 63     // 正常Cmd的Mask
)

type PacketHead struct {
	Cmd   uint8
	Param uint8
	Size  uint16
}

type PacketRead struct {
	PacketHead
	Comp    bool // 是否压缩
	Cont    bool // 是否连续
	HasHead bool // 是否已读
	Buf     *bytes.Buffer
}

func (x *PacketRead) ParseCmd() {
	x.Comp = (x.Cmd & PacketHeadCompFlag) > 0
	x.Cont = (x.Cmd & PacketHeadContFlag) > 0
	x.Cmd = x.Cmd & PacketHeadCmdMask
}

type Processor struct {
	buf  bytes.Buffer
	read PacketRead
}

func (this *Processor) Marshal(i interface{}) ([]byte, error) {
	if msg, ok := i.(Cmder); ok {
		msg.Init()
		bts, err := json.Marshal(msg)
		if err != nil {
			log.Printf("序列化化错误,%v", err)
			return nil, err
		}

		c := msg.GetCmd()
		if len(bts) > CompCmdLen {
			/*c |= PacketHeadCompFlag
			bts = glib.Compress(bts)*/
		}

		buf := new(bytes.Buffer)
		for i := 1; i > 0; {
			tmp := bts
			var cont uint8
			if len(bts) >= MaxCmdLen {
				tmp = bts[:MaxCmdLen]
				bts = bts[MaxCmdLen:]
				cont = PacketHeadContFlag
			} else {
				i = 0
			}
			BinWrite(buf, c|cont)
			BinWrite(buf, msg.GetParam())
			BinWrite(buf, uint16(len(tmp)))
			buf.Write(tmp)
		}
		return buf.Bytes(), nil
	} else {
		return nil, errors.New("no gcmd Cmder")
	}
}

func (this *Processor) UnMarshal(p []byte, mq chan interface{}) (n int, err error) {
	this.buf.Write(p)
	for {
		if !this.read.HasHead {
			if this.buf.Len() < PacketHeadLen {
				break
			}

			BinRead(&this.buf, &this.read.Cmd)
			BinRead(&this.buf, &this.read.Param)
			BinRead(&this.buf, &this.read.Size)
			this.read.HasHead = true
			this.read.ParseCmd()
		}
		if this.buf.Len() < int(this.read.Size) {
			break
		}
		this.read.HasHead = false

		data := this.buf.Next(int(this.read.Size))
		if this.read.Cont {
			if this.read.Buf == nil {
				this.read.Buf = &bytes.Buffer{}
			}
			this.read.Buf.Write(data)
			continue
		} else if this.read.Buf != nil && this.read.Buf.Len() > 0 {
			this.read.Buf.Write(data)
			data = this.read.Buf.Bytes()
			this.read.Buf = nil
		}
		if this.read.Comp {
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

func BinRead(buf *bytes.Buffer, data interface{}) {
	binary.Read(buf, binary.LittleEndian, data)
}

func BinWrite(buf *bytes.Buffer, data interface{}) {
	binary.Write(buf, binary.LittleEndian, data)
}

func NewProcessor() *Processor {
	return &Processor{}
}
