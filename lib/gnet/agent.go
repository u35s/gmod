package gnet

import (
	"bytes"
	"encoding/binary"
	"errors"
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
	Conn    net.Conn
	Process Processor

	lock          sync.Mutex
	unpackBuf     bytes.Buffer
	read          packetRead
	once          sync.Once
	errHandler    func(error)
	reciveChannel chan interface{}
	sendChannel   chan interface{}
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
	for send := range this.sendChannel {
		if send == nil {
			return errors.New("send is nil")
		}
		bts, err := this.Process.Marshal(send)
		if err != nil {
			return err
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
		this.unpack(bts[:num], this.reciveChannel)
	}
}

func (this *Agent) handleError(err error, end func()) {
	this.once.Do(func() {
		if this.errHandler != nil {
			this.errHandler(err)
		}
		if end != nil {
			end()
		}
	})
}

func (this *Agent) run() {
	go func() {
		this.handleError(this.recive(), func() {
			this.sendChannel <- nil
		})

	}()
	go func() {
		this.handleError(this.send(), func() {
			this.Conn.Close()
		})
	}()
}

func (this *Agent) Close(err error) {
	this.handleError(err, func() {
		this.sendChannel <- nil
		this.Conn.Close()
	})
}

func (this *Agent) GetMsg() <-chan interface{} {
	return this.reciveChannel
	/*select {
	case itfc := <-this.reciveChannel:
		ch := make(chan interface{}, 1)
		ch <- itfc
		return ch
	}*/
}

func (this *Agent) SendCmd(m interface{}) {
	this.sendChannel <- m
}

func (this *Agent) SetReciveChannel(ch chan interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for {
		select {
		case msg := <-this.reciveChannel:
			ch <- msg
		default:
			this.reciveChannel = ch
			return
		}
	}
}

func NewAgent(conn net.Conn, process Processor, errHandler func(error)) *Agent {
	g := &Agent{
		Conn: conn,

		Process:       process,
		errHandler:    errHandler,
		reciveChannel: make(chan interface{}, 1<<10),
		sendChannel:   make(chan interface{}, 1<<10),
	}
	g.run()
	return g
}
