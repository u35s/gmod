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

func GetString(idx string) string {
	return GetGroup(defaultGroup).GetString(idx)
}

func GetInt(idx string) int {
	return GetGroup(defaultGroup).GetInt(idx)
}

func GetUint(idx string) uint {
	return GetGroup(defaultGroup).GetUint(idx)
}

func GetFloat(idx string) float32 {
	return GetGroup(defaultGroup).GetFloat(idx)
}

func GetGroup(idx string) group {
	if group, ok := conf.groups[idx]; ok {
		return group
	}
	return make(group)
}
