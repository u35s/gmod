package server

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"

	"github.com/u35s/gmod"
	"github.com/u35s/gmod/examples/game/testcmd"
	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/lib/gnet"
	"github.com/u35s/gmod/mods/gsrvm"
)

type sessionServer struct {
	gmod.ModBase

	serverMsgChannel chan interface{}
}

func (this *sessionServer) Init() {
	defaultServerRoute()
	this.serverMsgChannel = make(chan interface{}, 1<<16)
	gsrvm.AddToListenAddr(":8001")
}

func (this *sessionServer) Wait() bool {
	gsrvm.EachToListenAddr(func(s *gsrvm.ToListenAddr) {
		if !s.Ok && this.listenTo(s.Addr) == nil {
			s.Ok = true
		}
	})
	return true
}

func (this *sessionServer) listenTo(addr string) error {
	listener, err := gnet.Listen(addr)
	if err != nil {
		log.Printf("listen %v err:%v", addr, err)
		return err
	}
	go gnet.Accept(listener, this.handleConn)
	return nil
}

func (this *sessionServer) handleConn(conn net.Conn) {
	srv := &gsrvm.ConnectedServer{}
	srv.Agent = gnet.NewAgent(conn, gcmd.NewProcessor(), func(err error) {
		gsrvm.Remove(srv)
		log.Printf("server %v,%v remote addr %v error %v",
			srv.Type, srv.Name, srv.Agent.Conn.RemoteAddr(), err)
	})

	select {
	case itfc := <-srv.Agent.GetMsg():
		if msg, ok := itfc.(*gcmd.CmdMessage); ok {
			var rev testcmd.CmdServer_establishConnection
			json.Unmarshal(msg.Data, &rev)
			srv.Type, srv.Name = rev.Type, rev.Name
			srv.Agent.SetReciveChannel(this.serverMsgChannel)
			gsrvm.Add(srv)
		}
	case <-time.After(2 * time.Second):
		err := errors.New("connection verify time out")
		log.Printf("%v", err)
		srv.Agent.Close(err)
	}
}

func (this *sessionServer) Run() {
	this.dealServerMsg()
}

func (this *sessionServer) dealServerMsg() {
	for {
		select {
		case msg := <-this.serverMsgChannel:
			deliverServerMsg(msg)
		default:
			return
		}
	}
}
