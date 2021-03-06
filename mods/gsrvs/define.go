package gsrvs

import (
	"sync"

	"github.com/u35s/gmod/lib/gnet"
)

type uint = uint64

type ServerBase struct {
	ID   uint
	Type string
	Name string
}

type ConnectedServer struct {
	ServerBase
	Agent *gnet.Agent
}

type ConnectedServerSlc []*ConnectedServer

type ToListenAddr struct {
	Net  string
	Addr string
	Ok   bool
}

type ToConnectServer struct {
	ServerBase
	Net       string
	Addr      string
	LocalType string
	LocalName string
	Ok        bool
}

type ConnectedServerSlcSafeMap struct {
	lock sync.RWMutex

	toListen  []*ToListenAddr
	toConnect []*ToConnectServer

	connectedServerSlcMap map[string]ConnectedServerSlc
	size                  int
}

func (this *ConnectedServerSlcSafeMap) init() {
	this.connectedServerSlcMap = make(map[string]ConnectedServerSlc, 0)
}
