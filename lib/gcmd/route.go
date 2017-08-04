package gcmd

var delivers [255][255]func(*CmdMessage, ...interface{})

func Route(cmd Cmder, h func(*CmdMessage, ...interface{})) {
	cmd.Init()
	delivers[cmd.GetCmd()][cmd.GetParam()] = h
}

func DeliverMsg(itfc ...interface{}) {
	if msg, ok := itfc[0].(*CmdMessage); ok {
		if h := delivers[msg.GetCmd()][msg.GetParam()]; h != nil {
			h(msg, itfc[1:]...)
		}
	}
}
