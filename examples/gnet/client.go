package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/u35s/gmod/examples/gnet/testcmd"
	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/lib/gnet"
)

func main() {
	conn, err := net.DialTimeout("tcp", ":8100", 2*time.Second)
	if err != nil {
		log.Print(err)
		return
	}
	handleConn(conn)

}
func handleConn(conn net.Conn) {
	log.Printf("connection server success,local addr %v, remote addr %v", conn.LocalAddr(), conn.RemoteAddr())
	run := true
	agent := gnet.NewAgent(conn, gcmd.NewProcessor(), func(itfc interface{}) {
		if msg, ok := itfc.(*gcmd.CmdMessage); ok {
			var rev testcmd.CmdServer_chat
			json.Unmarshal(msg.Data, &rev)
			fmt.Printf("server say %v\n", rev.Cnt)
		}
	}, func(err error) {
		run = false
		log.Printf("agent error,%v\n", err)
	})
	var send testcmd.CmdServer_chat
	send.Cnt = "hello"
	agent.SendCmd(&send)
	for run {
		time.Sleep(time.Second)
	}
}
