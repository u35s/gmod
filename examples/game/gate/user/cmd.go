package user

import (
	"encoding/json"
	"log"

	"github.com/u35s/gmod/examples/game/testcmd"
	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/mods/gsrvs"
)

var delivers [255][255]func(*user, *gcmd.CmdMessage)

func route(cmd gcmd.Cmder, h func(*user, *gcmd.CmdMessage)) {
	cmd.Init()
	delivers[cmd.GetCmd()][cmd.GetParam()] = h
}

func deliver(u *user, msg *gcmd.CmdMessage) {
	if h := delivers[msg.GetCmd()][msg.GetParam()]; h != nil {
		h(u, msg)
	}
}

func defaultRoute() {
	route(&testcmd.CmdUser_login{}, func(u *user, msg *gcmd.CmdMessage) {
		var rev testcmd.CmdUser_login
		json.Unmarshal(msg.Data, &rev)
		u.accid = rev.Accid
		var send testcmd.CmdServer_userLogin
		send.Accid = rev.Accid
		send.Seqid = u.seqid
		log.Printf("user login,accid %v,seqid %v", rev.Accid, send.Seqid)
		sendCmdToSession(&send)
	})
}

func sendCmdToSession(send gcmd.Cmder) {
	gsrvs.SendCmdToServer("session", "session", &send)
}
