package log

// Debugf 格式化日志
func Debugf(s string, v ...interface{}) {
	Logger.Debugf(s, v...)
}

// Infof 格式化日志
func Infof(s string, v ...interface{}) {
	Logger.Infof(s, v...)
}

// Warnf 格式化日志
func Warnf(s string, v ...interface{}) {
	Logger.Warnf(s, v...)
}

// Errorf 格式化日志
func Errorf(s string, v ...interface{}) {
	Logger.Errorf(s, v...)
}

// Fatalf 格式化日志
func Fatalf(s string, v ...interface{}) {
	Logger.Panicf(s, v...)
}

// Panicf 格式化日志
func Panicf(s string, v ...interface{}) {
	Logger.Panicf(s, v...)
}

// Debug 打日志
func Debug(v ...interface{}) {
	Logger.Debug(v...)
}

// Info 打日志
func Info(v ...interface{}) {
	Logger.Info(v...)
}

// Warn 打日志
func Warn(v ...interface{}) {
	Logger.Warn(v...)
}

// Error 打日志
func Error(v ...interface{}) {
	Logger.Error(v...)
}

// Panic 打日志
func Panic(v ...interface{}) {
	Logger.Panic(v...)
}
