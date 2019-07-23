# 日志库

 * 日志滚动，保存15天内的日志，每1小时(整点)分割一次日志（可配置）
 * 日志文件分级保存 debug, info, error
 * 可选是否输出日志到控制台
 * 支持输出打日志的文件文件行号

## 使用

只支持单例log

```go

// 使用前必须先初始化
log.Init(&log.LogConfig{
	Path:    ".runtime/logs", // 日志保存路径
	Name:    "api", // 日志文件名
	Console: true, // 是否把日志输出到控制台
	DebugConsole: true, // 是否把调试日志也输出到控制台
	Caller: true, // 是否输出打日志的文件和行号，会影响性能
	MaxAge: time.Hour * 24 * 15, // 保存多久的日志，默认15天
	RotationTime: time.Hour, // 多久分割一次日志，默认一小时
})

log.Debug("is debug log")
log.Info("is info log")
log.Warn("is warn log")
log.Error("is error log")
log.Panic("is panic log")

log.Debugf("is debug log %s %d %v", "string", 123, map[string]string{"test": "hello"})
log.Infof("is info log %s %d %v", "string", 123, map[string]string{"test": "hello"})
log.Warnf("is warn log %s %d %v", "string", 123, map[string]string{"test": "hello"})
log.Errorf("is error log %s %d %v", "string", 123, map[string]string{"test": "hello"})
log.Fatalf("is fatal log %s %d %v", "string", 123, map[string]string{"test": "hello"})
log.Panicf("is panic log %s %d %v", "string", 123, map[string]string{"test": "hello"})

```

## godoc

[API 文档](https://gowalker.org/github.com/go-eyas/eyas/log)