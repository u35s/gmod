package server

import (
	"log"
	"net"
	"time"

	"github.com/u35s/gmod"
	"github.com/u35s/gmod/examples/game/testcmd"
	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/lib/gnet"
	"github.com/u35s/gmod/mods/gsrvs"
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
	gsrvs.AddToConnectServer("session", "session", ":8001", "tcp")
}

func (this *gateServer) Wait() bool {
	gsrvs.EachToConnectServer(func(s *gsrvs.ToConnectServer) {
		if !s.Ok && this.connectTo(s.Addr, s.Type, s.Name) == nil {
			s.Ok = true
		}
	})
	return gsrvs.ServerSize() == 1
}

func (this *gateServer) connectTo(addr, tp, name string) error {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		log.Printf("connect to %v err:%v", addr, err)
		return err
	}
	srv := &gsrvs.ConnectedServer{ServerBase: gsrvs.ServerBase{Type: tp, Name: name}}
	srv.Agent = gnet.NewAgent(conn, gcmd.NewProcessor(), func(itfc interface{}) {
		this.serverMsgChannel <- itfc
	}, func(err error) {
		gsrvs.Remove(srv)
		log.Printf("server %v,%v remote addr %v error %v",
			srv.Type, srv.Name, srv.Agent.Conn.RemoteAddr(), err)
	})
	var send testcmd.CmdServer_establishConnection
	send.Type = this.tp
	send.Name = this.name
	srv.Agent.SendCmd(&send)
	gsrvs.Add(srv)
	return nil
}

func (this *gateServer) Run() {
	this.dealServerMsg()
}

func (this *gateServer) dealServerMsg() {
	for {
		select {
		case msg := <-this.serverMsgChannel:
			gcmd.DeliverMsg(msg)
		default:
			return
		}
	}
}
