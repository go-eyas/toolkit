# websocket 服务

开箱即用的 websocket 服务

## 协议

#### 心跳

心跳包为长度为 0 的空数据包，最长时间 30s 发一次，否则链接将会被断开

#### 请求响应数据

请求和响应的数据必须按照该协议来

**请求数据**

```json
{
  "cmd": "register",
  "seqno": "unique string",
  "data": {}
}
```

 * cmd 命令名称
 * seqno 请求标识符
 * data 请求数据


**响应数据**

```json
{
  "cmd": "register",
  "seqno": "unique string",
  "msg": "ok",
  "status": 0,
  "data": {}
}
```

 * cmd 命令名称，原样返回
 * seqno 请求标识符，原样返回
 * msg 处理后的消息，如果消息是处理成功的，默认都是 ok
 * status 错误状态码，0 为成功，非 0 为失败
 * data 响应数据

## 使用

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
  server.UseRequest(func(c *wsrv.Context) {
    if c.CMD != "register" {
      _, ok := c.Get("uid").(int)
      if !ok {
        c.Abort()
      }
    }
  })

  server.Handle("register", func(c *wsrv.Context) {
    c.Set("uid", 1001)
    
    // server push
    for sid, vals := range server.Session {
      if uid, ok := vals["uid"]; ok {
        server.Push(sid, &wsrv.WSResponse{
          CMD: "have_user_register",
          Data: map[string]interface{}{
            "uid": uid,
          },
        })
      } 
    }

    // server push current connection
    c.Push(&wsrv.WSResponse{
      CMD: "user_register",
      Data: map[string]interface{}{
        "uid": uid,
      },
    })

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

## API

[API 文档](https://gowalker.org/github.com/go-eyas/toolkit/websocket/wsrv)