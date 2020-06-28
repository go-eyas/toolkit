# websocket

# 使用

```go
import (
  "net/http"
  "github.com/go-eyas/toolkit/websocket"
)
func main() {
  ws := websocket.New(&Config{
    MsgType: websocket.TextMessage, // 消息类型 websocket.TextMessage | websocke.BinaryMessage
  })

  http.HandleFunc("/ws", ws.HTTPHandler)

  go func() {
    rec := ws.Receive()
    for {
      req, _ := <-rec
      req.Response([]byte("1234556"))
    }
  }()

  http.ListenAndServe("127.0.0.1:8800", nil)
}
```

# 服务

已经准备了一个开箱即用的服务，该服务按照特定协议工作，[详情请查看](./wsrv)

示例概览

```go
import (
  "net/http"
  "github.com/go-eyas/toolkit/websocket"
  "github.com/go-eyas/toolkit/websocket/wsrv"
)
func main() {
  server := wsrv.New(&Config{
    MsgType: websocket.TextMessage, // 消息类型 websocket.TextMessage | websocke.BinaryMessage
  })
  server.Use(func(c *wsrv.Context) {
    if c.CMD != "register" {
      _, ok := c.Get("uid").(int)
      if !ok {
        c.Abort()
      }
    }
  })

  server.Handle("register", func(c *wsrv.Context) {
    c.Set("uid", 1001)
    c.OK()
  })
  server,Handle("userinfo", func(c *wsrv.Context) {
    uid := c.Get("uid").(int)
    c.OK(GetUserInfoByID(uid))
  })

  http.HandleFunc("/ws", server.Engine.HTTPHandler)
  http.ListenAndServe("127.0.0.1:8800", nil)
}
```

## 协议



## godoc

[API 文档](https://gowalker.org/github.com/go-eyas/toolkit/websocket)