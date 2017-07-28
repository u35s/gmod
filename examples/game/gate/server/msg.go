package server

import (
	"encoding/json"
	"log"

	"github.com/u35s/gmod/examples/game/gate/userm"
	"github.com/u35s/gmod/examples/game/testcmd"
	"github.com/u35s/gmod/lib/gcmd"
)

var serverDelivers [255][255]func(*gcmd.CmdMessage)

func serverRoute(cmd gcmd.Cmder, h func(*gcmd.CmdMessage)) {
	cmd.Init()
	serverDelivers[cmd.GetCmd()][cmd.GetParam()] = h
}

func deliverServerMsg(itfc interface{}) {
	if msg, ok := itfc.(*gcmd.CmdMessage); ok {
		if h := serverDelivers[msg.GetCmd()][msg.GetParam()]; h != nil {
			h(msg)
		}
	}
}

func defaultServerRoute() {
	serverRoute(&testcmd.CmdServer_userLogin{}, func(msg *gcmd.CmdMessage) {
		var rev testcmd.CmdServer_userLogin
		json.Unmarshal(msg.Data, &rev)
		u := userm.GetVerifyUserBySeqid(rev.Seqid)
		if u != nil {
			userm.AddUser(u)
			log.Printf("user login success, accid %v,seqid %v,", rev.Accid, rev.Seqid)
			var send testcmd.CmdUser_chat
			send.Cnt = "login success"
			u.SendCmdToMe(&send)
		}
	})
}
