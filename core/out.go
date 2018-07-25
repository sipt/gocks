package core

import "fmt"

var Notify INotify

type INotify interface {
	Error(...interface{})
	Info(...interface{})
}

var Logger ILogger = &stdLogger{}

type ILogger interface {
	Trace(...interface{})
	Debug(...interface{})
	Info(...interface{})
	Error(...interface{})
}

type stdLogger struct{}

func (s *stdLogger) Trace(params ...interface{}) {
	fmt.Println("[TRACE]", fmt.Sprint(params ...))
}
func (s *stdLogger) Debug(params ...interface{}) {
	fmt.Println("[DEBUG]", fmt.Sprint(params ...))
}
func (s *stdLogger) Info(params ...interface{}) {
	fmt.Println("[INFO]", fmt.Sprint(params ...))
}
func (s *stdLogger) Error(params ...interface{}) {
	fmt.Println("[ERROR]", fmt.Sprint(params ...))
}
