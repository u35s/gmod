package gtime

import "time"

const (
	NanosecondN  uint = 1                   //纳秒
	MicrosecondN      = 1000 * NanosecondN  //微妙
	MillisecondN      = 1000 * MicrosecondN //毫秒
	SecondN           = 1000 * MillisecondN
	MinuteN           = 60 * SecondN

	SecondS = 1
	MinuteS = 60 * SecondS
	HourS   = 60 * MinuteS
	DayS    = 24 * HourS
	WeekS   = 7 * DayS
	MonthS  = 30 * DayS
)

func Time() uint {
	return uint(time.Now().Unix())
}
func TimeNano() uint {
	return uint(time.Now().UnixNano())
}
