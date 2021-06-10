package lorm

import (
	"fmt"
	"log"
	"os"
)

var (
	Log = Logger{}
)

func init() {

	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.SetFlags(log.LstdFlags | log.Llongfile)
	Log = Logger{log: logger}
}

type Logger struct {
	log *log.Logger
}

func (l *Logger) Fatalln(msg string, v ...interface{}) {
	arr := make([]interface{}, 0)
	arr = append(arr, msg)
	arr = append(arr, "\n")
	arr = append(arr, v...)
	l.log.Output(2, fmt.Sprintln(arr...))
	os.Exit(1)
}

func (l *Logger) Panicln(msg string, v ...interface{}) {
	arr := make([]interface{}, 0)
	arr = append(arr, msg)
	arr = append(arr, "\n")
	arr = append(arr, v...)
	l.log.Output(2, fmt.Sprintln(arr...))
	panic(msg)
}

func (l *Logger) Println(msg string, v ...interface{}) {


	arr := make([]interface{}, 0)
	arr = append(arr, msg)
	arr = append(arr, "\n")
	arr = append(arr, v...)
	l.log.Output(2, fmt.Sprintln(arr...))
}
