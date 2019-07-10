# 日志库

 * 日志滚动，保存15天内的日志，每1小时(整点)分割一次日志（可配置）
 * 日志文件分级保存 debug, info, error

## 使用

只支持单例log

```go

// 使用前必须先初始化
log.Init(&log.LogConfig{
	Path:    ".runtime/logs",
	Name:    "api",
	Console: true,
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