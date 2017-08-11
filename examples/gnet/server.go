package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/u35s/gmod/examples/gnet/testcmd"
	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/lib/gnet"
)

func main() {
	listener, err := gnet.Listen(":8100")
	if err != nil {
		log.Print(err)
		return
	}
	gnet.Accept(listener, handleConn)
}

func handleConn(conn net.Conn) {
	log.Printf("receive new connection,local addr %v,remote addr %v", conn.LocalAddr(), conn.RemoteAddr())
	var agent *gnet.Agent
	agent = gnet.NewAgent(conn, gcmd.NewProcessor(), func(itfc interface{}) {
		if msg, ok := itfc.(*gcmd.CmdMessage); ok {
			var rev testcmd.CmdServer_chat
			json.Unmarshal(msg.Data, &rev)
			agent.SendCmd(&rev)
		}
	}, func(err error) {
		log.Printf("agent error,%v\n", err)
	})
	var send testcmd.CmdServer_chat
	send.Cnt = "hello"
	agent.SendCmd(&send)
}
