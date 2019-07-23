# Rabbitmq 封装

封装 amqp 协议的基本使用方法，让amqp用起来更简单

## 使用

### 配置项 

```go
// Config 配置项
// ExchangeName 和 Exchange 二选一，用于指定发布和订阅时使用的交换机
type Config struct {
	Addr string // rabbitmq 地址
	ExchangeName string    // 使用该值创建一个直连的交换机
	Exchange     *Exchange // 自定义默认交换机
	Consumer *Consumer // 在定于队列时，作为消费者使用的参数
}
```

### 交换机

```go
const (
    ExchangeDirect  = "direct"
    ExchangeFanout  = "fanout"
    ExchangeTopic   = "topic"
    ExchangeHeaders = "headers"
)

type Exchange struct {
	Name       string // 名称
	Kind       string // 交换机类型，4 种类型之一
	Durable    bool   // 是否持久化
	AutoDelete bool   // 是否自动删除
	Internal   bool   // 是否内置,如果设置 为true,则表示是内置的交换器,客户端程序无法直接发送消息到这个交换器中,只能通过交换器路由到交换器的方式
	NoWait     bool   // 是否等待通知定义交换机结果
	Args       amqp.Table
}
```

### 队列

```go
type Queue struct {
	Name       string     // 必须包含前缀标识使用类型 msg. | rpc. | reply. | notify.
	Durable    bool       // 消息代理重启后，队列依旧存在
	AutoDelete bool       // 当最后一个消费者退订后即被删除
	Exclusive  bool       // 只被一个连接（connection）使用，而且当连接关闭后队列即被删除
	NoWait     bool       // 不需要服务器返回
	ReplyTo    *Queue     // rpc 的消息回应道哪个队列
	Args       amqp.Table // 一些消息代理用他来完成类似与TTL的某些额外功能
}
```

### 消息结构

```go
// Message 消息体
type Message struct {
	ContentType string // 消息类型
	Queue       *Queue // 来自于哪个队列
	Data        []byte // 消息数据
}
```

### 示例

```go
import (
	"github.com/go-eyas/eyas/amqp"
)

func main() {
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
}

```

## godoc

[API 文档](https://gowalker.org/github.com/go-eyas/eyas/amqp)
