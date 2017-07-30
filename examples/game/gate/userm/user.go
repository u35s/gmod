package userm

import (
	"errors"
	"log"

	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/lib/gnet"
	"github.com/u35s/gmod/lib/gtime"
)

type user struct {
	Agent *gnet.Agent
	Err   chan error

	loginTime uint

	seqid uint
	accid uint
}

func (this *user) SendCmdToMe(send gcmd.Cmder) {
	this.Agent.SendCmd(send)
}

func (this *user) verify() {
	this.refresh()

	if this.loginTime+gtime.MinuteS < gtime.Time() {
		this.Err <- errors.New("verify time out")
	}
}

func (this *user) deliverMsg() {
	for {
		select {
		case itfc := <-this.Agent.GetMsg():
			if msg, ok := itfc.(*gcmd.CmdMessage); ok {
				log.Printf("deliver user %v msg,cmd %v,param %v", this.seqid, msg.GetCmd(), msg.GetParam())
				deliver(this, msg)
			}
		case err := <-this.Err:
			this.destory(err)
			return
		default:
			return
		}
	}
}

func (this *user) refresh() {
	this.deliverMsg()
}

func (this *user) destory(err error) {
	userm.removeVerifyUser(this)
	userm.removeUser(this)
	log.Printf("user accid %v,seqid %v destory,err %v", this.accid, this.seqid, err)
}
