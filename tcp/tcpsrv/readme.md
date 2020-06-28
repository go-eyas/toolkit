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

    // 自定义tcp数据包协议，实现下方两个方法即可
    // 将业务数据封装成tcp数据包
    // Packer: func(data []byte) ([]byte,  error) {}, 
    // 将 tcp 连接收到的数据包解析成业务数据，返回的业务数据必须符合上方定义的 json 数据
    // Parser: func(conn *tcp.Conn, pack []byte) ( [][]byte,  error) {}, 
  })
  if err != nil {
    panic(err)
  }
  
  // log 中间件
  server.Use(func(c *tcpsrv.Context) {
    fmt.Printf("TCP 收到 cmd=%s seqno=%s data=%s\n", c.CMD, c.Seqno, string(c.Payload))
    c.Next()
    fmt.Printf("TCP 响应 cmd=%s seqno=%s data=%s\n", c.CMD, c.Seqno, string(c.Response.Data))
  })

  // 验证中间件
  server.Use(func(c *tcpsrv.Context) {
    if c.CMD != "register" {
      _, ok := c.Get("uid").(int64)
      if !ok {
        c.Response.Msg = "permission defined"
        c.Response.Status = 401
        c.Abort() // 停止后面的中间件执行
        return
      } 
    }
    c.Next() // 如后续无操作，可省略
  })

  server.Handle("register", func(c *tcpsrv.Context) {
    body := &struct {
      UID int64 `json:"uid"`
    }{}
    err := c.Bind(body) // 绑定json数据
    if err != nil {
      panic(err) // 在 Handle panic 后不会导致程序异常，会响应错误数据到客户端
    }
    c.Set("uid", body.UID) // 设置该连接的会话值
    c.OK()
  })

  server.Handle("userinfo", func(c *tcpsrv.Context) {
    uid := c.Get("uid").(int64) // 获取会话值
    c.OK(findUserByUID(uid)) // OK 可设置响应数据，如果不设置
  })
}
```

### 客户端

客户端是该协议的实现，在符合上述协议的服务器都可使用

```go
package main

import (
  "github.com/go-eyas/toolkit/tcp"
  "github.com/go-eyas/toolkit/tcp/tcpsrv"
  "fmt"
)

func main()  {
    client, err := tcpsrv.NewClientSrv(&tcp.Config{
      Addr:    ":6601",

      // 自定义tcp数据包协议，实现下方两个方法即可
      // 将业务数据封装成tcp数据包
      // Packer: func(data []byte) ([]byte,  error) {}, 
      // 将 tcp 连接收到的数据包解析成业务数据，返回的业务数据必须符合上方定义的 json 数据
      // Parser: func(conn *tcp.Conn, pack []byte) ( [][]byte,  error) {}, 
    })
    if err != nil {
      panic(err)
    }
  
    // 每当服务器发送了数据过来，都会以 cmd 作为时间名触发事件
    client.On("register", func(response *tcpsrv.TCPResponse) {
      fmt.Println("on receive register msg:", response)
    })
  
    client.On("userinfo", func(response *tcpsrv.TCPResponse) {
      fmt.Println("on receive userinfo msg:", response)
    })
  
    // send 发送后，会等待服务器的响应，res 为服务器的响应数据
    res, err := client.Send("register", map[string]interface{}{
      "uid": 1234,
    })
    if err != nil {
      panic(err)
    }
    fmt.Println("send register response: ", res)
  
    res, err = client.Send("userinfo")
    if err != nil {
      panic(err)
    }
    // 响可直接解析绑定 data 数据
    res.BindJSON(&struct {
     UID int64
    }{})
    fmt.Println("send userinfo response: ", res)
}
```

# API

[API 文档](https://gowalker.org/github.com/go-eyas/toolkit/tcp/tcpsrv)