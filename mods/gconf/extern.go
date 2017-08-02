package gconf

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strings"
)

var conf *config
var defaultGroup string = "default"
var lastGroup string = defaultGroup

func init() {
	conf = new(config)
	conf.groups = make(map[string]group)
	conf.groups[defaultGroup] = make(group)
}

func ReadFile(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	buf.Write(data)
	reader := bufio.NewReader(buf)
	for {
		bts, _, err := reader.ReadLine()
		if err != nil {
			return err
		}
		if len(bts) == 0 || bts[0] == '#' {
			continue
		}
		str := strings.TrimSpace(string(bts))
		if gk, gv, ok := groupAction(str); ok {
			lastGroup = gk
			conf.groups[lastGroup] = gv
		} else if optk, optv, ok := optionAction(str); ok {
			conf.groups[lastGroup][optk] = optv
		}
	}
	return nil
}

func String(idx string) (ret string) {
	eachGroupBreak(func(g *group) bool {
		if ret = g.GetString(idx); len(ret) > 0 {
			return true
		}
		return false
	})
	return
}

func Int(idx string) (ret int) {
	eachGroupBreak(func(g *group) bool {
		if ret = g.GetInt(idx); ret > 0 {
			return true
		}
		return false
	})
	return
}

func Uint(idx string) (ret uint) {
	eachGroupBreak(func(g *group) bool {
		if ret = g.GetUint(idx); ret > 0 {
			return true
		}
		return false
	})
	return
}

func Float(idx string) (ret float32) {
	eachGroupBreak(func(g *group) bool {
		if ret = g.GetFloat(idx); ret > 0 {
			return true
		}
		return false
	})
	return
}

func eachGroupBreak(brk func(*group) bool) {
	for _, group := range conf.groups {
		if brk(&group) {
			break
		}
	}
}

func Group(idx string) group {
	if group, ok := conf.groups[idx]; ok {
		return group
	}
	return make(group)
}
