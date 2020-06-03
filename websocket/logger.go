package websocket

type loggerI interface {
	Log(...interface{})
	Logf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}

type l struct{}

func (l) Log(v ...interface{})              {}
func (l) Logf(s string, v ...interface{})   {}
func (l) Error(v ...interface{})            {}
func (l) Errorf(s string, v ...interface{}) {}

var emptyLogger = &l{}

var logger loggerI = emptyLogger
