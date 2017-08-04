package server

import (
	"encoding/json"
	"log"

	"github.com/u35s/gmod/examples/game/gate/user"
	"github.com/u35s/gmod/examples/game/testcmd"
	"github.com/u35s/gmod/lib/gcmd"
)

func serverRoute(cmd gcmd.Cmder, f func(*gcmd.CmdMessage)) {
	gcmd.Route(cmd, func(h func(*gcmd.CmdMessage)) func(*gcmd.CmdMessage, ...interface{}) {
		return func(msg *gcmd.CmdMessage, itfc ...interface{}) {
			h(msg)
		}
	}(f))
}

func defaultServerRoute() {
	serverRoute(&testcmd.CmdServer_userLogin{}, func(msg *gcmd.CmdMessage) {
		var rev testcmd.CmdServer_userLogin
		json.Unmarshal(msg.Data, &rev)
		u := user.GetVerifyUserBySeqid(rev.Seqid)
		if u != nil {
			user.AddUser(u)
			log.Printf("user login success, accid %v,seqid %v,", rev.Accid, rev.Seqid)
			var send testcmd.CmdUser_chat
			send.Cnt = "login success"
			u.SendCmdToMe(&send)
		}
	})
}
