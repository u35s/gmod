package gsrvm

import (
	"net"

	"github.com/u35s/gmod"
	"github.com/u35s/gmod/lib/gnet"
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

func SetNewAgentFunc(h func(net.Conn) *gnet.Agent) {
	srvm.newAgentFunc = h
}

func SendCmdToServer(tp, name string, msg interface{}) {
	srvm.sendCmdToServer(tp, name, msg)
}
