package gsrvm

import (
	"log"

	"github.com/u35s/gmod/lib/gnet"
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

type connectedServer struct {
	Type  string
	Name  string
	Agent *gnet.Agent
	Add   bool
}

type connectedServerSlc []*connectedServer

func (this *connectedServerSlc) Add(srv *connectedServer) {
	*this = append(*this, srv)
	log.Printf("%v,%v add success,local addr %v,remote addr %v",
		srv.Type, srv.Name, srv.Agent.Conn.LocalAddr(), srv.Agent.Conn.RemoteAddr())
}

func (this *connectedServerSlc) Remove(srv *connectedServer) {
	for i := range *this {
		*this = append((*this)[:i], (*this)[i+1:]...)
		log.Printf("%v,%v remove success,local addr %v,remote addr %v",
			srv.Type, srv.Name, srv.Agent.Conn.LocalAddr(), srv.Agent.Conn.RemoteAddr())
	}
}

type connectedServerSlcMap map[string]connectedServerSlc
