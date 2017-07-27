package gsrvm

import (
	"github.com/u35s/gmod"
)

var srvm *serverManager

func Mod() gmod.Moder {
	if srvm == nil {
		srvm = new(serverManager)
	}
	return srvm
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
