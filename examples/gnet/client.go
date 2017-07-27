package main

import (
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
	agent := gnet.NewAgent(conn, gcmd.NewProcessor())
	var send testcmd.CmdServer_chat
	send.Cnt = "hello"
	agent.SendChannel <- &send
	for {
		select {
		case v := <-agent.ReciveChannel:
			if msg, ok := v.(*gcmd.CmdMessage); ok {
				var rev testcmd.CmdServer_chat
				agent.Processor.Unmarshal(msg.Data, &rev)
				fmt.Printf("server say %v\n", rev.Cnt)
			}
		case err := <-agent.Err:
			log.Printf("agent error,%v\n", err)
		}
	}
}
