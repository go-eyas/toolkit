package websocket

type LoggerI interface {
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}

type l struct{}

func (l) Info(v ...interface{})             {}
func (l) Infof(s string, v ...interface{})  {}
func (l) Error(v ...interface{})            {}
func (l) Errorf(s string, v ...interface{}) {}

var EmptyLogger = &l{}

var logger LoggerI = EmptyLogger
