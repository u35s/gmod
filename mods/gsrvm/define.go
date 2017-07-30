package gsrvm

import (
	"sync"

	"github.com/u35s/gmod"
	"github.com/u35s/gmod/lib/gnet"
)

type ConnectedServer struct {
	Type  string
	Name  string
	Agent *gnet.Agent
}

type ConnectedServerSlc []*ConnectedServer

type ToListenAddr struct {
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

type ConnectedServerSlcSafeMap struct {
	gmod.ModBase
	lock sync.RWMutex

	toListen  []*ToListenAddr
	toConnect []*ToConnectServer

	connectedServerSlcMap map[string]ConnectedServerSlc
	size                  int
}

func (this *ConnectedServerSlcSafeMap) Init() {
	this.connectedServerSlcMap = make(map[string]ConnectedServerSlc, 0)
}
