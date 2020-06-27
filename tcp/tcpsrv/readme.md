# TCP 服务

开箱即用的 TCP 服务

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

# 使用

示例概览

### 服务器

```go
package main

import (
  "github.com/go-eyas/toolkit/tcp"
  "github.com/go-eyas/toolkit/tcp/tcpsrv"
  "fmt"
)

func main() {
  server, err := tcpsrv.NewServerSrv(&tcp.Config{
    Network:":6700",
  })
  if err != nil {
    panic(err)
  }
  
  // log 中间件
  server.Use(func(c *tcpsrv.Context) {
    fmt.Printf("TCP 收到 cmd=%s seqno=%s data=%s", c.CMD, c.Seqno, string(c.Payload))
    c.Next()
    fmt.Printf("TCP 响应 cmd=%s seqno=%s data=%s", c.CMD, c.Seqno, string(c.Response.Data))
  })

  // 验证中间件
  server.Use(func(c *tcpsrv.Context) {
    if c.CMD != "register" {
      _, ok := c.Get("uid").(int64)
      if !ok {
        c.Response.Status = 401
        c.Abort()
        return
      } 
    }
    c.Next() // 如后续无操作，可省略
  })

  server.Handle("register", func(c *tcpsrv.Context) {
    c.Set("uid", int64(100001))
    c.OK()
  })

  server.Handle("userinfo", func(c *tcpsrv.Context) {
    uid := c.Get("uid").(int64)
    c.OK(findUserByUID(uid))
  })
}
```