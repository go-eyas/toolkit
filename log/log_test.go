package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	Init(&LogConfig{
		Path:         ".runtime/logs",
		Name:         "api",
		Console:      true,
		Caller:       true,
		DebugConsole: true,
	})

	Debug("is debug log")
	Info("is info log")
	Warn("is warn log")
	Error("is error log")
	// Panic("is panic log")

	Debugf("is debug log %s %d %v", "string", 123, map[string]string{"test": "hello"})
	Infof("is info log %s %d %v", "string", 123, map[string]string{"test": "hello"})
	Warnf("is warn log %s %d %v", "string", 123, map[string]string{"test": "hello"})
	Errorf("is error log %s %d %v", "string", 123, map[string]string{"test": "hello"})
	// Fatalf("is fatal log %s %d %v", "string", 123, map[string]string{"test": "hello"})
	// Panicf("is panic log %s %d %v", "string", 123, map[string]string{"test": "hello"})
}
