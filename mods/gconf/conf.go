package gconf

import (
	"strings"

	"github.com/u35s/gmod/lib/utils"
)

type config struct {
	groups map[string]group
}

type group map[string]string

func (g group) GetString(idx string) string {
	return g[idx]
}

func (g group) GetInt(idx string) int {
	return utils.Atoi(g.GetString(idx))
}

func (g group) GetUint(idx string) uint {
	return utils.Atou(g.GetString(idx))
}

func (g group) GetFloat(idx string) float32 {
	return utils.Atof(g.GetString(idx))
}

func groupAction(line string) (string, group, bool) {
	length := len(line)
	line = strings.TrimPrefix(line, "[")
	line = strings.TrimSuffix(line, "]")
	if len(line)+2 != length {
		return "", nil, false
	}
	group := make(group)
	return line, group, true
}

func optionAction(line string) (string, string, bool) {
	strs := strings.Split(line, "=")
	if len(strs) != 2 {
		return "", "", false
	}
	return strings.TrimSpace(strs[0]), strings.TrimSpace(strs[1]), true
}
