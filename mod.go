package gmod

type Moder interface {
	Init()
	Wait() bool
	Run()
	End()
}

type ModBase struct{}

func (this *ModBase) Init()      {}
func (this *ModBase) Wait() bool { return true }
func (this *ModBase) Run()       {}
func (this *ModBase) End()       {}
