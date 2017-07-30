package server

import (
	"log"
	"net"
	"time"

	"github.com/u35s/gmod"
	"github.com/u35s/gmod/examples/game/testcmd"
	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/lib/gnet"
	"github.com/u35s/gmod/mods/gsrvm"
)

type gateServer struct {
	gmod.ModBase
	tp               string
	name             string
	serverMsgChannel chan interface{}
}

func (this *gateServer) Init() {
	defaultServerRoute()
	this.tp, this.name = "gate", "gate"
	this.serverMsgChannel = make(chan interface{}, 1<<16)
	gsrvm.AddToConnectServer("session", "session", ":8001")
}

func (this *gateServer) Wait() bool {
	gsrvm.EachToConnectServer(func(s *gsrvm.ToConnectServer) {
		if !s.Ok && this.connectTo(s.Addr, s.Type, s.Name) == nil {
			s.Ok = true
		}
	})
	return gsrvm.ServerSize() == 1
}

func (this *gateServer) connectTo(addr, tp, name string) error {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		log.Printf("connect to %v err:%v", addr, err)
		return err
	}
	srv := &gsrvm.ConnectedServer{Type: tp, Name: name}
	srv.Agent = gnet.NewAgent(conn, gcmd.NewProcessor(), func(err error) {
		gsrvm.Remove(srv)
		log.Printf("server %v,%v remote addr %v error %v",
			srv.Type, srv.Name, srv.Agent.Conn.RemoteAddr(), err)
	})
	srv.Agent.SetReciveChannel(this.serverMsgChannel)
	var send testcmd.CmdServer_establishConnection
	send.Type = this.tp
	send.Name = this.name
	srv.Agent.SendCmd(&send)
	gsrvm.Add(srv)
	return nil
}

func (this *gateServer) Run() {
	this.dealServerMsg()
}

func (this *gateServer) dealServerMsg() {
	for {
		select {
		case msg := <-this.serverMsgChannel:
			deliverServerMsg(msg)
		default:
			return
		}
	}
}
