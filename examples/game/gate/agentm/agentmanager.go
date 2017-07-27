package agentm

import (
	"log"
	"net"

	"github.com/u35s/gmod"
	"github.com/u35s/gmod/examples/game/gate/userm"
	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/lib/gnet"
)

type agentManager struct {
	gmod.ModBase
}

func (this *agentManager) Init() {
	listener, err := gnet.Listen(":8002")
	if err != nil {
		log.Print(err)
		return
	}
	go gnet.Accept(listener, this.handleConn)
}

func (this *agentManager) handleConn(conn net.Conn) {
	log.Printf("[agentm],receive new connection,local addr %v,remote addr %v", conn.LocalAddr(), conn.RemoteAddr())
	agent := gnet.NewAgent(conn, gcmd.NewProcessor())
	u := userm.NewVerifyUser(agent)
	userm.AddVerifyUser(u)
}
