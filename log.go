package lorm

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Logger struct {
	log *log.Logger
}

func (l *Logger) Fatalln(msg string, v ...any) {
	arr := make([]any, 0)
	arr = append(arr, msg)
	arr = append(arr, "\n")
	arr = append(arr, v...)
	l.log.Output(2, fmt.Sprintln(arr...))
	os.Exit(1)
}

func (l *Logger) Panicln(msg string, v ...any) {
	arr := make([]any, 0)
	arr = append(arr, msg)
	arr = append(arr, "\n")
	arr = append(arr, v...)
	l.log.Output(2, fmt.Sprintln(arr...))
	panic(msg)
}

func (l *Logger) Println(msg string, v ...any) {
	if !strings.HasPrefix(msg, "\n") {
		msg = "\n" + msg
	}
	arr := make([]any, 0)
	arr = append(arr, msg)
	arr = append(arr, "\n")
	arr = append(arr, v...)

	l.log.Output(2, fmt.Sprintln(arr...))
}

// todo 下面未重构--------------
