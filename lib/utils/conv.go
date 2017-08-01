package utils

import "strconv"

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
