package utils

import (
	"strconv"
	"strings"
)

type uint = uint64

func Atoi(s string) int {
	if i, err := strconv.ParseInt(s, 10, 0); err == nil {
		return int(i)
	}
	return 0
}

func Atou(s string) uint {
	if i, err := strconv.ParseUint(s, 10, 0); err == nil {
		return uint(i)
	}
	return 0
}

func Atof(s string) float32 {
	if f, err := strconv.ParseFloat(s, 32); err == nil {
		return float32(f)
	}
	return 0
}

func Split(str string, seps string) []string {
	if len(seps) == 0 {
		return []string{str}
	}
	sepSlc := make([]string, len(seps))
	for i := 0; i < len(seps); i++ {
		sepSlc[i] = string(seps[i])
	}
	return mulit_split([]string{str}, sepSlc)
}

func mulit_split(strSlc []string, sepSlc []string) []string {
	if len(sepSlc) == 0 {
		return strSlc
	}

	subStr := make([]string, 0, 16)
	sep := sepSlc[0]
	for i := 0; i < len(strSlc); i++ {
		ret := strings.Split(strSlc[i], sep)
		for j := range ret {
			if len(ret[j]) > 0 {
				subStr = append(subStr, ret[j])
			}
		}
	}
	return mulit_split(subStr, sepSlc[1:])
}
