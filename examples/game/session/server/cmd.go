package server

import (
	"encoding/json"
	"log"

	"github.com/u35s/gmod/examples/game/testcmd"
	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/mods/gsrvs"
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
		log.Printf("user login verify success,accid %v,seqid %v", rev.Accid, rev.Seqid)
		gsrvs.SendCmdToServer("gate", "gate", &rev)
	})
	serverRoute(&testcmd.CmdServer_forwardUserMsg{}, func(msg *gcmd.CmdMessage) {
		var rev testcmd.CmdServer_forwardUserMsg
		json.Unmarshal(msg.Data, &rev)
		subMsg := new(gcmd.CmdMessage)
		subMsg.SetBase(rev.SubCmd, rev.SubParam)
		subMsg.Data = rev.Data
		deliverUserMsg(rev.UID, subMsg)
	})
}

func deliverUserMsg(uid uint, msg *gcmd.CmdMessage) {
	if msg.GetCmd() == testcmd.CmdUser {
		switch msg.GetParam() {
		}
	}
}
