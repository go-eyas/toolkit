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
		Logger.Debugf("%s %s", getCaller(), fmt.Sprintf(s, v...))
	} else {
		Logger.Debugf(s, v...)
	}
}

// Infof 格式化日志
func Infof(s string, v ...interface{}) {
	if printCaller {
		Logger.Infof("%s %s", getCaller(), fmt.Sprintf(s, v...))
	} else {
		Logger.Infof(s, v...)
	}
}

// Warnf 格式化日志
func Warnf(s string, v ...interface{}) {
	if printCaller {
		Logger.Warnf("%s %s", getCaller(), fmt.Sprintf(s, v...))
	} else {
		Logger.Warnf(s, v...)
	}
}

// Errorf 格式化日志
func Errorf(s string, v ...interface{}) {
	if printCaller {
		Logger.Errorf("%s %s", getCaller(), fmt.Sprintf(s, v...))
	} else {
		Logger.Errorf(s, v...)
	}
}

// Fatalf 格式化日志
func Fatalf(s string, v ...interface{}) {
	if printCaller {
		Logger.Fatalf("%s %s", getCaller(), fmt.Sprintf(s, v...))
	} else {
		Logger.Fatalf(s, v...)
	}
}

// Panicf 格式化日志
func Panicf(s string, v ...interface{}) {
	if printCaller {
		Logger.Panicf("%s %s", getCaller(), fmt.Sprintf(s, v...))
	} else {
		Logger.Panicf(s, v...)
	}
}

// Debug 打日志
func Debug(v ...interface{}) {
	if printCaller {
		Logger.Debugf("%s %s", getCaller(), fmt.Sprint(v...))
	} else {
		Logger.Debug(v...)
	}
}

// Info 打日志
func Info(v ...interface{}) {
	if printCaller {
		Logger.Infof("%s %s", getCaller(), fmt.Sprint(v...))
	} else {
		Logger.Info(v...)
	}
}

// Warn 打日志
func Warn(v ...interface{}) {
	if printCaller {
		Logger.Warnf("%s %s", getCaller(), fmt.Sprint(v...))
	} else {
		Logger.Warn(v...)
	}
}

// Error 打日志
func Error(v ...interface{}) {
	if printCaller {
		Logger.Errorf("%s %s", getCaller(), fmt.Sprint(v...))
	} else {
		Logger.Error(v...)
	}
}

// Panic 打日志
func Panic(v ...interface{}) {
	if printCaller {
		Logger.Panicf("%s %s", getCaller(), fmt.Sprint(v...))
	} else {
		Logger.Panic(v...)
	}
}
