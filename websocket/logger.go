package websocket

type loggerI interface {
	Log(...interface{})
	Logf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

type l struct{}

func (l) Log(v ...interface{})              {}
func (l) Logf(s string, v ...interface{})   {}
func (l) Fatal(v ...interface{})            {}
func (l) Fatalf(s string, v ...interface{}) {}

var emptyLogger = &l{}

var logger loggerI = emptyLogger
