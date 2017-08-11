package utils

import "log"

var logSeverity = LogAll

type logLevel uint

const (
	LogAll logLevel = iota
	LogErr
	LogWrn
	LogInf
	LogDbg
	LogMax
)

func logf(severity logLevel, format string, args ...interface{}) {
	if severity >= logSeverity {
		log.Printf(format, args...)
	}
}

func SetServerity(lvl logLevel) {
	if lvl < LogMax {
		logSeverity = lvl
	}
}

func Dbg(format string, args ...interface{}) {
	logf(LogDbg, format, args...)
}

func Inf(format string, args ...interface{}) {
	logf(LogInf, format, args...)
}

func Wrn(format string, args ...interface{}) {
	logf(LogWrn, format, args...)
}

func Err(format string, args ...interface{}) {
	logf(LogErr, format, args...)
}
