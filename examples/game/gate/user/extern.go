package user

import (
	"github.com/u35s/gmod"
	"github.com/u35s/gmod/lib/gtime"
)

var seqid uint = 0
var userm *userManager

func NewVerifyUser() *user {
	seqid++
	return &user{seqid: seqid, loginTime: gtime.Time(), Err: make(chan error, 2)}
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

func Mod() gmod.Moder {
	if userm == nil {
		userm = new(userManager)
	}
	return userm
}
