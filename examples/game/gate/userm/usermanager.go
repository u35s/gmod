package userm

import (
	"sync"

	"github.com/u35s/gmod"
)

type userManager struct {
	gmod.ModBase

	rwlock    sync.RWMutex
	verifings map[uint]*user
	accids    map[uint]*user
}

func (this *userManager) Init() {
	this.verifings = make(map[uint]*user)
	this.accids = make(map[uint]*user)
}

func (this *userManager) Run() {
	for _, u := range this.verifings {
		u.verify()
	}
	for _, u := range this.accids {
		u.refresh()
	}
}

func (this *userManager) addUser(u *user) {
	this.removeVerifyUser(u)
	this.rwlock.Lock()
	this.accids[u.accid] = u
	this.rwlock.Unlock()
}

func (this *userManager) removeUser(u *user) {
	this.rwlock.Lock()
	delete(this.accids, u.accid)
	this.rwlock.Unlock()
}

func (this *userManager) addVerifyUser(u *user) {
	this.rwlock.Lock()
	this.verifings[u.seqid] = u
	this.rwlock.Unlock()
}

func (this *userManager) getUserBySeqid(id uint) (*user, bool) {
	this.rwlock.RLock()
	defer this.rwlock.RUnlock()
	u, ok := this.verifings[id]
	return u, ok
}

func (this *userManager) removeVerifyUser(u *user) {
	this.rwlock.Lock()
	delete(this.verifings, u.seqid)
	this.rwlock.Unlock()
}
