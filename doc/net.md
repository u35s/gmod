# gnet
  gnet提供了对建立tcp连接的简单封装及tcp消息包的收发,其中Processor接口允许你自定义消息包，此接口中Marshal的作用是把消息序列化为二进制类型，然后通过对应的tcp连接发送出去,相对的tcp连接调用Unmarshal提取对应的消息来完成服务器之间或者玩家与服务器之间的消息交换.

  gnet是完全独立于gmod的，这意味着gnet可以用在任何通过tcp交换信息的地方，例如一个简单的client<->server通信的例子(相应的代码在examples/gnet 文件夹下)：

client
```
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/u35s/gmod/examples/gnet/testcmd"
	"github.com/u35s/gmod/gcmd"
	"github.com/u35s/gmod/gnet"
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
				json.Unmarshal(msg.Data, &rev)
				fmt.Printf("server say %v\n", rev.Cnt)
			}
		case err := <-agent.Err:
			log.Printf("agent error,%v\n", err)
		}
	}
}
```
server
```
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
```
