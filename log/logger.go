package log

import (
	"fmt"
	"runtime"
)

func getCaller() string {
	if !printCaller {
		return ""
	}
	funcName, _, line, ok := runtime.Caller(2)
	caller := ""
	if ok {
		caller = fmt.Sprintf("%s:%d", runtime.FuncForPC(funcName).Name(), line)
	}

	return caller
}

// Debugf 格式化日志
func Debugf(s string, v ...interface{}) {
	if printCaller {
		SugaredLogger.Debugf("%s %s", getCaller(), fmt.Sprintf(s, v...))
	} else {
		SugaredLogger.Debugf(s, v...)
	}
}

// Infof 格式化日志
func Infof(s string, v ...interface{}) {
	if printCaller {
		SugaredLogger.Infof("%s %s", getCaller(), fmt.Sprintf(s, v...))
	} else {
		SugaredLogger.Infof(s, v...)
	}
}

// Warnf 格式化日志
func Warnf(s string, v ...interface{}) {
	if printCaller {
		SugaredLogger.Warnf("%s %s", getCaller(), fmt.Sprintf(s, v...))
	} else {
		SugaredLogger.Warnf(s, v...)
	}
}

// Errorf 格式化日志
func Errorf(s string, v ...interface{}) {
	if printCaller {
		SugaredLogger.Errorf("%s %s", getCaller(), fmt.Sprintf(s, v...))
	} else {
		SugaredLogger.Errorf(s, v...)
	}
}

// Fatalf 格式化日志
func Fatalf(s string, v ...interface{}) {
	if printCaller {
		SugaredLogger.Fatalf("%s %s", getCaller(), fmt.Sprintf(s, v...))
	} else {
		SugaredLogger.Fatalf(s, v...)
	}
}

// Panicf 格式化日志
func Panicf(s string, v ...interface{}) {
	if printCaller {
		SugaredLogger.Panicf("%s %s", getCaller(), fmt.Sprintf(s, v...))
	} else {
		SugaredLogger.Panicf(s, v...)
	}
}

// Debug 打日志
func Debug(v ...interface{}) {
	if printCaller {
		SugaredLogger.Debugf("%s %s", getCaller(), fmt.Sprint(v...))
	} else {
		SugaredLogger.Debug(v...)
	}
}

// Info 打日志
func Info(v ...interface{}) {
	if printCaller {
		SugaredLogger.Infof("%s %s", getCaller(), fmt.Sprint(v...))
	} else {
		SugaredLogger.Info(v...)
	}
}

// Warn 打日志
func Warn(v ...interface{}) {
	if printCaller {
		SugaredLogger.Warnf("%s %s", getCaller(), fmt.Sprint(v...))
	} else {
		SugaredLogger.Warn(v...)
	}
}

// Error 打日志
func Error(v ...interface{}) {
	if printCaller {
		SugaredLogger.Errorf("%s %s", getCaller(), fmt.Sprint(v...))
	} else {
		SugaredLogger.Error(v...)
	}
}

// Panic 打日志
func Panic(v ...interface{}) {
	if printCaller {
		SugaredLogger.Panicf("%s %s", getCaller(), fmt.Sprint(v...))
	} else {
		SugaredLogger.Panic(v...)
	}
}
