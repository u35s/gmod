package gmod

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/u35s/gmod/lib/gtime"
)

const frameWarningTime = 100 * gtime.MillisecondN
const baseFrame uint = 40

type serverState uint

func (this *serverState) String() string {
	switch *this {
	case serverStateInit:
		return "init"
	case serverStateWait:
		return "wait"
	case serverStateRun:
		return "runing"
	case serverStateEnd:
		return "stopping"
	}
	return "null"
}

const (
	ServerStateNull serverState = iota
	serverStateInit
	serverStateWait
	serverStateRun
	serverStateEnd
)

type server struct {
	mods         []Moder
	state        serverState
	frameRunTime uint
}

func (this *server) start() {
	this.setFrameRunTime(baseFrame * 2)
}

func (this *server) init() {
	var timer gtime.Timer
	timer.Reset()
	for i := range this.mods {
		this.mods[i].Init()
	}
	log.Printf("server mod init,time %v second", timer.Elapse()/gtime.SecondN)
}

func (this *server) wait() bool {
	for i := range this.mods {
		if !this.mods[i].Wait() {
			return false
		}
	}
	return true
}

func (this *server) run() {
	for i := range this.mods {
		this.mods[i].Run()
	}
}

func (this *server) end() {
	var timer gtime.Timer
	timer.Reset()
	for i := range this.mods {
		this.mods[i].End()
	}
	log.Printf("server mod end,time %v second", timer.Elapse()/gtime.SecondN)
}

func (this *server) signal() {
	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT)
	select {
	case <-sigChan:
		this.setState(serverStateEnd)
	}
}

func (this *server) loop() {
	var timer gtime.Timer
	for {
		timer.Reset()
		switch this.state {
		case serverStateInit:
			this.init()
			this.setState(serverStateWait)
		case serverStateWait:
			if this.wait() {
				this.setState(serverStateRun)
			}
		case serverStateRun:
			this.run()
		case serverStateEnd:
			this.end()
			return
		default:
		}
		e := timer.Elapse()
		if e < this.frameRunTime {
			sleep := this.frameRunTime - e
			time.Sleep(time.Duration(sleep))
		} else if e > frameWarningTime {
			log.Printf("this frame run time in %v millisecond", e/gtime.MillisecondN)
		}

	}
}

func (this *server) setState(state serverState) {
	this.state = state
	log.Printf("set server sate %v", state.String())
}

func (this *server) setFrameRunTime(runTime uint) {
	this.frameRunTime = runTime * gtime.MillisecondN
}

var srv *server

func init() {
	srv = new(server)
	srv.start()
	go srv.signal()
}

func Run(mods ...Moder) {
	srv.mods = append(srv.mods, mods...)
	srv.setState(serverStateInit)
	srv.loop()
}
