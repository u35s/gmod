package server

import (
	"encoding/json"
	"log"

	"github.com/u35s/gmod/examples/game/testcmd"
	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/mods/gsrvs"
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

func deliverUserMsg(uid uint, msg *gcmd.CmdMessage) {
	if msg.GetCmd() == testcmd.CmdUser {
		switch msg.GetParam() {
		}
	}

}

func defaultServerRoute() {
	serverRoute(&testcmd.CmdServer_forwardUserMsg{}, func(msg *gcmd.CmdMessage) {
		var rev testcmd.CmdServer_forwardUserMsg
		json.Unmarshal(msg.Data, &rev)
		subMsg := new(gcmd.CmdMessage)
		subMsg.SetBase(rev.SubCmd, rev.SubParam)
		subMsg.Data = rev.Data
		deliverUserMsg(rev.UID, subMsg)
	})
	serverRoute(&testcmd.CmdServer_userLogin{}, func(msg *gcmd.CmdMessage) {
		var rev testcmd.CmdServer_userLogin
		json.Unmarshal(msg.Data, &rev)
		log.Printf("user login verify success,accid %v,seqid %v", rev.Accid, rev.Seqid)
		var send testcmd.CmdServer_userLogin
		send.Seqid = rev.Seqid
		send.Accid = rev.Accid
		sendMsgToGate(&send)
	})
}

func sendMsgToGate(send gcmd.Cmder) {
	gsrvs.SendCmdToServer("gate", "gate", send)
}
