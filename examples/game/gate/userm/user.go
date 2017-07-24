package userm

import (
	"errors"
	"log"

	"github.com/u35s/gmod/gcmd"
	"github.com/u35s/gmod/gnet"
	"github.com/u35s/gmod/gtime"
)

type user struct {
	Agent *gnet.Agent

	activeClose error
	agentClose  error
	loginTime   uint

	seqid uint
	accid uint
}

func (this *user) verify() {
	if this.loginTime+gtime.MinuteS < gtime.Time() {
		this.activeClose = errors.New("verify time out")
	}
	this.refresh()
}

func (this *user) deliverMsg() {
	for {
		select {
		case v := <-this.Agent.ReciveChannel:
			if msg, ok := v.(*gcmd.CmdMessage); ok {
				log.Printf("deliver user %v msg,cmd %v,param %v", this.seqid, msg.GetCmd(), msg.GetParam())
				deliver(this, msg)
			}
		case err := <-this.Agent.Err:
			this.agentClose = err
			return
		default:
			return
		}
	}
}

func (this *user) refresh() {
	if this.isClose() {
		return
	}
	this.deliverMsg()
}

func (this *user) isClose() bool {
	if this.agentClose != nil || this.activeClose != nil {
		if this.activeClose != nil {
			this.Agent.Close(this.activeClose)
		}
		this.destory()
		return true
	}
	return false
}

func (this *user) destory() {
	userm.removeUser(this)
	log.Printf("user accid %v,seqid %v destory", this.accid, this.seqid)
}
