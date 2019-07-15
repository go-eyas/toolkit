# gin 工具箱

## 中间件 

#### error

捕获到在http处理时的错误

在 handler 和其他地方如果产生了 error 可直接panic，到这里统一处理，简化 if err != nil 之类的代码

```go
import (
  "github.com/go-eyas/toolkit/gin/middleware"
  "github.com/go-eyas/toolkit/log"
)

log.Init(&log.Config{})

// engine 
route.Use(middleware.Error(log.SugaredLogger)) // 如果不需要日志记录，传 nil

// handler
func HelloHandler(c *gin.Context) {
  panic("text") // {msg: "text", code: 0, data: {}}
  panic(gin.H{"code": 0, "msg": "some error"}) // {与传入的数据一致，} code 默认999999，status 默认 400，msg 默认 unknow error
  panic(errors.New("some error")) // {msg: "some error", code: 999999, data: {}}
  panic(Struct{...}) // {msg: "unknow", code: 999999, data: {...struct 数据}}
}
```

#### logger 

使用 zap 打印日志

log.Init(&log.Config{})

```go
import (
  "github.com/go-eyas/toolkit/gin/middleware"
  "github.com/go-eyas/toolkit/log"
)

log.Init(&log.Config{})

route.Use(middleware.Ginzap(log.SugaredLogger))
```
