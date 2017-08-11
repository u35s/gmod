package main

import (
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/u35s/gmod/examples/game/testcmd"
	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/lib/gnet"
)

func main() {
	conn, err := net.DialTimeout("tcp", ":8002", 2*time.Second)
	if err != nil {
		log.Print(err)
		return
	}
	handleConn(conn)

}
func handleConn(conn net.Conn) {
	log.Printf("connection gate success,local addr %v, remote addr %v", conn.LocalAddr(), conn.RemoteAddr())
	agent := gnet.NewAgent(conn, gcmd.NewProcessor(), func(itfc interface{}) {
		if msg, ok := itfc.(*gcmd.CmdMessage); ok {
			deliverMsg(msg)
		}
	}, func(err error) {
		log.Printf("agent error,%v\n", err)
	})
	var send testcmd.CmdUser_login
	send.Accid = 1
	agent.SendCmd(&send)
}

func deliverMsg(msg *gcmd.CmdMessage) {
	if msg.GetCmd() == testcmd.CmdUser {
		switch msg.GetParam() {
		case testcmd.CmdUserParam_chat:
			var rev testcmd.CmdUser_chat
			json.Unmarshal(msg.Data, &rev)
			log.Printf("server say %v", rev.Cnt)
		}
	}
}
