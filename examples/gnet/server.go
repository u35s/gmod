package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/u35s/gmod/examples/gnet/testcmd"
	"github.com/u35s/gmod/glib/gcmd"
	"github.com/u35s/gmod/glib/gnet"
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
	agent := gnet.NewAgent(conn, gcmd.NewProcessor())
	var send testcmd.CmdServer_chat
	send.Cnt = "welcome"
	agent.SendChannel <- &send
	for {
		select {
		case v := <-agent.ReciveChannel:
			if msg, ok := v.(*gcmd.CmdMessage); ok {
				var rev testcmd.CmdServer_chat
				json.Unmarshal(msg.Data, &rev)
				agent.SendChannel <- &rev
			}
		case err := <-agent.Err:
			log.Printf("agent err,%v", err)
		}
	}
}
