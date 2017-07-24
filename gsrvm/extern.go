package gsrvm

import (
	"github.com/u35s/gmod"
)

var Mod gmod.Moder

var srvm *serverManager

func init() {
	srvm = new(serverManager)
	Mod = srvm
}

func AddToListen(s ToListenServer) {
	srvm.toListen = append(srvm.toListen, s)
}

func AddToConnect(s ToConnectServer) {
	srvm.toConnect = append(srvm.toConnect, s)
}

func SetDealMsgFunc(h func(interface{})) {
	srvm.dealMsgFunc = h
}

func SendCmdToServer(tp, name string, msg interface{}) {
	srvm.sendCmdToServer(tp, name, msg)
}
