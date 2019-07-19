# go 工具箱

为了快速使用通用功能，做一次通用封装

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

err := redis.Init(&redis.Config{
  Cluster:  false, // 是否集群
  Addrs:    []string{"10.0.3.252:6379"}, // redis 地址，如果是集群则在数组上写多个元素
  Password: "",
	DB:       1,
})
redis.Set("tookit:test", `{"hello": "world"}`)
res, err := redis.Get("toolkit:test")
var data interface{}
res.JSON(&data)
```

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

# [Gin 中间件 & 工具](./gin)

```go
import "github.com/go-eyas/toolkit/gin/util" // 工具函数
import "github.com/go-eyas/toolkit/gin/middleware" // 中间件
```

# [工具函数 util](./util)

```go
import "github.com/go-eyas/toolkit/util"
```
