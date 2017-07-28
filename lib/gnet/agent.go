package gnet

import (
	"bytes"
	"encoding/binary"
	"net"
	"sync"
)

const (
	maxPacketLen   = 0x7fff // 最大包长度
	packetHeadLen  = 4
	packetContFlag = 1 << 15 //包是否连续
)

type packetRead struct {
	size    uint16
	cont    bool
	hasHead bool

	buf *bytes.Buffer
}

func (x *packetRead) parse() {
	x.cont = (x.size & packetContFlag) > 0
	x.size = x.size & maxPacketLen
}

type Agent struct {
	lock sync.Mutex

	unpackBuf bytes.Buffer
	read      packetRead

	Process       Processor
	Conn          net.Conn
	Err           chan error
	ReciveChannel chan interface{}
	SendChannel   chan interface{}
}

func binRead(buf *bytes.Buffer, data interface{}) {
	binary.Read(buf, binary.LittleEndian, data)
}

func binWrite(buf *bytes.Buffer, data interface{}) {
	binary.Write(buf, binary.LittleEndian, data)
}

func (this *Agent) pack(bts []byte) []byte {
	buf := bytes.Buffer{}
	for {
		if len(bts) > maxPacketLen {
			binWrite(&buf, uint16(maxPacketLen|packetContFlag))
			binWrite(&buf, bts[:maxPacketLen])
			bts = bts[maxPacketLen:]
		} else {
			binWrite(&buf, uint16(len(bts)))
			binWrite(&buf, bts)
			break
		}
	}
	return buf.Bytes()
}

func (this *Agent) unpack(bts []byte, ch chan interface{}) {
	this.unpackBuf.Write(bts)
	for {
		if !this.read.hasHead {
			if this.unpackBuf.Len() < packetHeadLen {
				break
			}
			binRead(&this.unpackBuf, &this.read.size)
			this.read.hasHead = true
			this.read.parse()

		}
		if this.unpackBuf.Len() < int(this.read.size) {
			break
		}
		this.read.hasHead = false

		data := this.unpackBuf.Next(int(this.read.size))
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
		itfc, err := this.Process.Unmarshal(data)
		if err == nil {
			ch <- itfc
		}
	}
}

func (this *Agent) send() error {
	for cmd := range this.SendChannel {
		if cmd == nil {
			return nil
		}
		bts, err := this.Process.Marshal(cmd)
		if err == nil {
		}
		bts = this.pack(bts)
		_, err = this.Conn.Write(bts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Agent) recive() error {
	bts := make([]byte, 1<<16)
	for {
		num, err := this.Conn.Read(bts)
		if err != nil {
			return err
		}
		this.unpack(bts[:num], this.ReciveChannel)
	}
}

func (this *Agent) run() {
	go func() {
		this.Err <- this.recive()
	}()
	go func() {
		this.Err <- this.send()
	}()
}

func (this *Agent) Close(err error) {
	this.SendChannel <- nil
	this.Conn.Close()
}

func (this *Agent) SendMsg(m interface{}) {
	this.SendChannel <- m
}

func (this *Agent) SetReciveChannel(ch chan interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for {
		select {
		case msg := <-this.ReciveChannel:
			ch <- msg
		default:
			this.ReciveChannel = ch
			return
		}
	}
}

func NewAgent(conn net.Conn, process Processor) *Agent {
	g := &Agent{
		Process:       process,
		Conn:          conn,
		Err:           make(chan error, 1<<4),
		ReciveChannel: make(chan interface{}, 1<<10),
		SendChannel:   make(chan interface{}, 1<<10),
	}
	g.run()
	return g
}
