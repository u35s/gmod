package gnet

import (
	"net"
	"sync"
)

type Agent struct {
	lock sync.Mutex

	Processor     Processor
	Conn          net.Conn
	Err           chan error
	ReciveChannel chan interface{}
	SendChannel   chan interface{}
}

func (this *Agent) send() error {
	for cmd := range this.SendChannel {
		if cmd == nil {
			return nil
		}
		bts, err := this.Processor.Pack(cmd)
		if err != nil {
			return err
		}
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
		_, err = this.Processor.UnPack(bts[:num], this.ReciveChannel)
		if err != nil {
			return err
		}
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

func NewAgent(conn net.Conn, p Processor) *Agent {
	g := &Agent{
		Conn:      conn,
		Processor: p,

		Err:           make(chan error, 1<<4),
		ReciveChannel: make(chan interface{}, 1<<10),
		SendChannel:   make(chan interface{}, 1<<10),
	}
	g.run()
	return g
}
