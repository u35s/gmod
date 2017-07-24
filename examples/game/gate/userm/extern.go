package userm

import (
	"github.com/u35s/gmod"
	"github.com/u35s/gmod/gnet"
	"github.com/u35s/gmod/gtime"
)

var seqid uint = 0
var userm *userManager

var Mod gmod.Moder

func NewVerifyUser(g *gnet.Agent) *user {
	seqid++
	return &user{Agent: g, seqid: seqid, loginTime: gtime.Time()}
}

func AddVerifyUser(u *user) {
	userm.addVerifyUser(u)
}

func GetVerifyUserBySeqid(id uint) *user {
	u, ok := userm.getUserBySeqid(id)
	if ok {
		return u
	}
	return nil
}

func AddUser(u *user) {
	userm.addUser(u)
}

func init() {
	userm = new(userManager)
	Mod = userm
	defaultRoute()
}
