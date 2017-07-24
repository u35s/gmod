package gsrvm

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"

	"github.com/u35s/gmod"
	"github.com/u35s/gmod/examples/game/testcmd"
	"github.com/u35s/gmod/gcmd"
	"github.com/u35s/gmod/gnet"
)

type serverManager struct {
	gmod.ModBase
	serverTypeCount  int
	addServerChannel chan *ConnectedServer
	serverMsgChannel chan interface{}
	connectedServers ConnectedServerSlcMap
	dealMsgFunc      func(interface{})
	toListen         []ToListenServer
	toConnect        []ToConnectServer
}

func (this *serverManager) Init() {
	this.addServerChannel = make(chan *ConnectedServer, 1<<4)
	this.serverMsgChannel = make(chan interface{}, 1<<16)
	this.connectedServers = make(ConnectedServerSlcMap)
}

func (this *serverManager) Wait() bool {
	for k, v := range this.toListen {
		if !v.Ok && this.listenTo(v.Addr) == nil {
			this.toListen[k].Ok = true
		}
	}
	for k, v := range this.toConnect {
		if !v.Ok && this.connectTo(v.Addr, v.Type, v.Name, v.LocalType, v.LocalName) == nil {
			this.toConnect[k].Ok = true
		}
	}
	this.addServer()
	return len(this.connectedServers) >= len(this.toConnect)
}

func (this *serverManager) Run() {
	this.addServer()
	this.dealServerMsg()
}

func (this *serverManager) sendCmdToServer(tp, name string, msg interface{}) {
	slc, ok := this.connectedServers[tp]
	if ok {
		for _, v := range slc {
			if name == "" || v.Name == name {
				v.Agent.SendChannel <- msg
			}
		}
	}
}

func (this *serverManager) dealServerMsg() {
	for {
		select {
		case msg := <-this.serverMsgChannel:
			if this.dealMsgFunc != nil {
				this.dealMsgFunc(msg)
			} else {
				log.Printf("srvm dealMsgFunc is nil")
			}
		default:
			return
		}
	}
}

func (this *serverManager) addServer() {
	for _, slc := range this.connectedServers {
		for _, v := range slc {
			select {
			case err := <-v.Agent.Err:
				log.Printf("server %v,%v remote addr %v error %v", v.Type, v.Name, v.Agent.Conn.RemoteAddr(), err)
				this.addServerChannel <- v
			default:
			}
		}
	}
	for {
		select {
		case srv := <-this.addServerChannel:
			slc, ok := this.connectedServers[srv.Type]
			if srv.Add {
				if !ok {
					slc = make(ConnectedServerSlc, 0, 1)
				}
				srv.Add = false
				slc.Add(srv)
			} else if ok {
				slc.Remove(srv)
			}
			this.connectedServers[srv.Type] = slc
		default:
			return
		}
	}
}

func (this *serverManager) connectTo(addr, tp, name, localTp, localName string) error {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		log.Printf("connect to %v err:%v", addr, err)
		return err
	}
	agent := gnet.NewAgent(conn, gcmd.NewProcessor())
	agent.SetReciveChannel(this.serverMsgChannel)
	var send testcmd.CmdServer_establishConnection
	send.Type = localTp
	send.Name = localName
	agent.SendChannel <- &send
	this.addServerChannel <- &ConnectedServer{Type: tp, Name: name, Agent: agent, Add: true}
	return nil
}

func (this *serverManager) listenTo(addr string) error {
	listener, err := gnet.Listen(addr)
	if err != nil {
		log.Printf("listen %v err:%v", addr, err)
		return err
	}
	go gnet.Accept(listener, this.handleConn)
	return nil
}

func (this *serverManager) handleConn(conn net.Conn) {
	agent := gnet.NewAgent(conn, gcmd.NewProcessor())
	select {
	case v := <-agent.ReciveChannel:
		if msg, ok := v.(*gcmd.CmdMessage); ok {
			var rev testcmd.CmdServer_establishConnection
			json.Unmarshal(msg.Data, &rev)
			agent.SetReciveChannel(this.serverMsgChannel)
			this.addServerChannel <- &ConnectedServer{Type: rev.Type, Name: rev.Name, Agent: agent, Add: true}
		}
	case err := <-agent.Err:
		log.Printf("handle connection err %v", err)
	case <-time.After(2 * time.Second):
		err := errors.New("connection verify time out")
		log.Printf("%v", err)
		agent.Close(err)
	}
}
