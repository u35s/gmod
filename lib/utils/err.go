package utils

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
)

func DumpStack(name string, id uint) {
	if err := recover(); err != nil {
		var buf bytes.Buffer
		bs := make([]byte, 1<<12)
		num := runtime.Stack(bs, false)
		buf.WriteString(fmt.Sprintf("Panic: %s\n", err))
		buf.Write(bs[:num])
		log.Print(buf.String())
	}
}

func PrintStack() {
	var buf bytes.Buffer
	bs := make([]byte, 1<<12)
	num := runtime.Stack(bs, false)
	buf.Write(bs[:num])
	log.Print(buf.String())
}
