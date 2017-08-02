package gcmd

var delivers [255][255]func(*CmdMessage)

func Route(cmd Cmder, h func(*CmdMessage)) {
	cmd.Init()
	delivers[cmd.GetCmd()][cmd.GetParam()] = h
}

func DeliverMsg(itfc interface{}) {
	if msg, ok := itfc.(*CmdMessage); ok {
		if h := delivers[msg.GetCmd()][msg.GetParam()]; h != nil {
			h(msg)
		}
	}
}
