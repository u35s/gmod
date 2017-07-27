package gsrvm

import (
	"log"

	"github.com/u35s/gmod/glib/gnet"
)

type ToListenServer struct {
	Addr string
	Ok   bool
}

type ToConnectServer struct {
	Type      string
	Name      string
	Addr      string
	LocalType string
	LocalName string
	Ok        bool
}

type ConnectedServer struct {
	Type  string
	Name  string
	Agent *gnet.Agent
	Add   bool
}

type ConnectedServerSlc []*ConnectedServer

func (this *ConnectedServerSlc) Add(srv *ConnectedServer) {
	*this = append(*this, srv)
	log.Printf("%v,%v add success,local addr %v,remote addr %v",
		srv.Type, srv.Name, srv.Agent.Conn.LocalAddr(), srv.Agent.Conn.RemoteAddr())
}

func (this *ConnectedServerSlc) Remove(srv *ConnectedServer) {
	for i := range *this {
		*this = append((*this)[:i], (*this)[i+1:]...)
		log.Printf("%v,%v remove success,local addr %v,remote addr %v",
			srv.Type, srv.Name, srv.Agent.Conn.LocalAddr(), srv.Agent.Conn.RemoteAddr())
	}
}

type ConnectedServerSlcMap map[string]ConnectedServerSlc
