# TCP 

基于底层的 tcp 封装，有 客户端和服务端，本工具的服务端和客户端封装了通用的操作，均可单独在用于各个项目。

只简化了常用的 tcp 操作

# 配置

实例化 tcp 的服务端和客户端的时候要传入配置，服务端和客户端配置共用同个结构体 `*tcp.Config` ，说明如下 

```go
// Config 配置项
type Config struct {
  Addr    string                                // tcp 地址，在客户端使用为需要连接的地址，在服务端使用为监听的地址
  Network string                                // tcp 的网络类型，可选值为 "tcp", "tcp4", "tcp6", "unix" or "unixpacket"
  Packer  func([]byte) ([]byte, error)          // tcp 数据包的封装函数，传入的数据是需要发送的业务数据，返回发送给 tcp 的数据
  Parser  func(*Conn, []byte) ([][]byte, error) // 将收到的数据包，根据私有协议转换成业务数据，在这里处理粘包,半包等数据包问题，返回处理好的数据包
}
```

Packer 和 Parser 是 tcp 数据包的处理，是同时对等出现的。用于实现 tcp 的私有协议。

如果 Packer 和 Parser 两个都不传，会使用默认的私有协议实现，默认私有协议的数据包组成是

```
header[4字节标识body字节长度] + body[任意长度]
```

# 客户端 Client

tcp 的客户端封装，带有自动重连，简易的 api 封装，只关注收发数据即可

```go
package main

import (
    "fmt"
	"github.com/go-eyas/toolkit/tcp"
)

func main() {
  client, err := tcp.NewClient(&tcp.Config{
    Network: "tcp", // 网络类型，不填默认 tcp
    // tcp 服务端地址
  	Addr:    "127.0.0.1:6600",
  
    // 私有协议实现，不传将使用默认的私有协议实现
  	// Parser: func([]byte) ([]byte, error) {},
  	// Packer: func(*Conn, []byte) ([][]byte, error){},
  })
  if err != nil {
  	panic(err)
  }
  
  // 接收数据
  ch := client.Receive()
  go func() {
    for msg := range ch {
      // msg.Data 经过 Parser 处理过的数据
      // msg.Conn tcp 连接实例 
      fmt.Println("client receive:", string(msg.Data))
    }
  }()
  
  // 发送数据，send 后将立马把数据传给 Packer 处理后，在发送到 tcp 连接
  err = client.Send([]byte("hello world1"))
}


```

# 服务端 Server

```go
package main

import (
    "fmt"
	"github.com/go-eyas/toolkit/tcp"
)

func main() {
  server, err := tcp.NewServer(&tcp.Config{
  	Network: "tcp", // 网络类型，不填默认 tcp
    // tcp 监听地址
    Addr:    "127.0.0.1:6600",
    
    // 私有协议实现，不传将使用默认的私有协议实现
    // Parser: func([]byte) ([]byte, error) {},
    // Packer: func(*Conn, []byte) ([][]byte, error){},
  })
  
  if err != nil {
  	panic(err)
  }
  
  // 接收数据
  ch := server.Receive()
  for data := range ch {
  	fmt.Printf("server receive: %v", data.Data)
    
    // 服务器收到数据后，响应发送一条数据到客户端 
  	err := data.Response([]byte("server receive your message"))
  	if err != nil {
  		panic(err)
  	}
  }

  // 给所有连接都发送消息
  for connID, conn := range server.Sockets {
    fmt.Println("connID: ", connID)
    server.Send(conn, []byte("broadcast some message"))
    // or
    // server.SendConnID(connID, []byte("broadcast some message"))
  }
}

```

# API

[API 文档](https://gowalker.org/github.com/go-eyas/toolkit/tcp)