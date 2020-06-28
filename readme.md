# go 工具箱

为了快速使用通用功能，做一次通用封装

# 使用


```
go get -u -v github.com/go-eyas/toolkit
```

# [HTTP 客户端 http](./http)

```go
import "github.com/go-eyas/toolkit/http"

github := http.BaseURL("https://api.github.com")
res, err := github.Get("/repos/eyasliu/blog/issues")
var data interface{}
res.JSON(&data)
```

# [日志 log](./log)

```go
import "github.com/go-eyas/toolkit/log"

log.Init(&log.Config{})
log.Info("log init ok")
log.Infof("is info log %s %d %v", "string", 123, map[string]string{"test": "hello"})
```

# [Redis](./redis)

```go
import "github.com/go-eyas/toolkit/redis"

r, err := redis.Init(&redis.Config{
  Cluster:  false, // 是否集群
  Addrs:    []string{"127.0.0.1:6379"}, // redis 地址，如果是集群则在数组上写多个元素
  Password: "",
	DB:       1,
})
r.Set("tookit:test", `{"hello": "world"}`)
str, err := r.Get("toolkit:test")

```

# [TCP](./tcp)

**server**

```go
import "github.com/go-eyas/toolkit/tcp"

func main() {
  server, err := tcp.NewServer(&tcp.Config{
  	Network: "tcp", // 网络类型，不填默认 tcp
    // tcp 监听地址
    Addr:    "127.0.0.1:6600",
    
    // 私有协议实现，不传将使用默认的私有协议实现
    // Parser: func([]byte) ([]byte, error) {},
    // Packer: func(*Conn, []byte) ([][]byte, error){},
  })
  
  // 接收数据
  ch := server.Receive()
  for data := range ch {
  	fmt.Printf("server receive: %v", data.Data)
    
    // 服务器收到数据后，响应发送一条数据到客户端 
  	err := data.Response([]byte("server receive your message"))
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

**client**

```go
import "github.com/go-eyas/toolkit/tcp"

func main() {
  client, err := tcp.NewClient(&tcp.Config{
    Network: "tcp", // 网络类型，不填默认 tcp
    // tcp 服务端地址
  	Addr:    "127.0.0.1:6600",
  
    // 私有协议实现，不传将使用默认的私有协议实现
  	// Parser: func([]byte) ([]byte, error) {},
  	// Packer: func(*Conn, []byte) ([][]byte, error){},
  })
  
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

还有一个开箱即用的 [tcp 服务](./tcp/tcpsrv)

# [长连接 Websocket](./websocket)

```go
import "github.com/go-eyas/toolkit/websocket"

ws := websocket.New(&Config{})
http.HandleFunc("/ws", ws.HTTPHandler)
go func() {
  rec := ws.Receive()
  for {
    req, _ := <-rec
    req.Response([]byte("1234556"))
  }
}()
http.ListenAndServe("127.0.0.1:8800", nil)
```

还有一个类似 http api 的[开箱即用服务](./websocket/wsrv)

# [RabbitMQ amqp](./amqp)

```go
import "github.com/go-eyas/toolkit/amqp"

mq := amqp.New(*amqp.Config{
	Addr: "amqp://guest:guest@127.0.0.1:5672",
	ExchangeName: "toolkit.exchange.test",
})
queue := &amqp.Queue{Name: "toolkit.queue.test"}
err := mq.Pub(queue, &amqp.Message{Data: []byte("{\"hello\":\"world\"}")})

msgch, err := mq.Sub(queue)
for msg := range msgch {
	fmt.Printf("%s", string(msg.Data))
}

```

# [配置项 config](./config)

```go
import "github.com/go-eyas/toolkit/config"

conf := struct {
  Host string
  Port int
}{}
config.Init("config", &conf)
```

# [数据库 ORM](./db)

```go
import "github.com/go-eyas/toolkit/db"

var db *gorm.DB = db.Gorm(&db.Config{"mysql", "username:password@127.0.0.1:3306/test"})
var db *xorm.Engine = db.Xorm(&db.Config{"mysql", "username:password@127.0.0.1:3306/test"})

defer db.Close()
```

# [资源模型CRUD](./db/resource)

资源自动 crud

```go
import "github.com/go-eyas/toolkit/db"
import "github.com/go-eyas/toolkit/db/resource"

type Article struct {
  ID      int64  `resource:"pk;search:none"`
  Title   string `resource:"create;update;search:like"`
  Status  byte   `resource:"search:="`
}

var r, db, err =  resource.New(&db.Config{"mysql", "username:password@127.0.0.1:3306/test"}, Article{})

r.Create(map[string]string{"title": "hello eyas"}) // 增
r.Delete(1) // 删
r.Update(1, map[string]int{"status": 1}) // 改

// 查，指定主键
var m = &Article{}
r.Detail(1, m)

// 查，指定查询条件查列表
var list = []*Article{}
r.List(&list, map[string]interface{}{"title": "he"}, []string{"id DESC"})


```

# [Gin 中间件 & 工具](./gin)

```go
import "github.com/go-eyas/toolkit/gin/util" // 工具函数
import "github.com/go-eyas/toolkit/gin/middleware" // 中间件
```

# [事件分发器 Emitter](./emit)

```go
import "github.com/go-eyas/toolkit/emit"
fn1 := func(data interface{}) {
  fmt.Printf("fn1 receive data: %v", data)
}

emit.On("evt", fn1).Off("evt", fn1)
emit.Emit("evt", "hello emitter")

```

# [邮件发送 Email](./email)

```go
import (
  "github.com/go-eyas/toolkit/email"
  "github.com/BurntSushi/toml"
)

func ExampleSample() {
	tomlConfig := `
host = "smtp.qq.com"
port = "465"
account = "893521870@qq.com"
password = "haha, wo cai bu gao su ni ne"
name = "unit test"
secure = true
[tpl.a]
bcc = ["Jeason <eyasliu@163.com>"] # 抄送
cc = [] # 抄送人
subject = "Welcome, {{.Name}}" # 主题
text = "Hello, I am {{.Name}}" # 文本
html = "<h1>Hello, I am {{.Name}}</h1>" # html 内容
`
	conf := &Config{}
	toml.Decode(tomlConfig, conf)
	email := New(conf)
	email.SendByTpl("Yuesong Liu <liuyuesongde@163.com>", "a", struct{ Name string }{"Batman"})
}
```

# [工具函数 util](./util)

```go
import "github.com/go-eyas/toolkit/util"
```
