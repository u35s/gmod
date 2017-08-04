package user

import (
	"encoding/json"
	"log"

	"github.com/u35s/gmod/examples/game/testcmd"
	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/mods/gsrvs"
)

func userRoute(cmd gcmd.Cmder, f func(*gcmd.CmdMessage, *user)) {
	gcmd.Route(cmd, func(h func(*gcmd.CmdMessage, *user)) func(*gcmd.CmdMessage, ...interface{}) {
		return func(msg *gcmd.CmdMessage, itfc ...interface{}) {
			if u, ok := itfc[0].(*user); ok {
				h(msg, u)
			}
		}
	}(f))
}

func defaultRoute() {
	userRoute(&testcmd.CmdUser_login{}, func(msg *gcmd.CmdMessage, u *user) {
		var rev testcmd.CmdUser_login
		json.Unmarshal(msg.Data, &rev)
		u.accid = rev.Accid
		var send testcmd.CmdServer_userLogin
		send.Accid = rev.Accid
		send.Seqid = u.seqid
		log.Printf("user login,accid %v,seqid %v", rev.Accid, send.Seqid)
		gsrvs.SendCmdToServer("session", "session", &send)
	})
}
