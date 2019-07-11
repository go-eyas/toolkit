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

## godoc

```
package websocket // import "github.com/go-eyas/toolkit/websocket"


CONSTANTS

const (
        // 文本消息
        TextMessage = 1

        // 二进制数据消息
        BinaryMessage = 2
)

TYPES

type Config struct {
        MsgType         int                      // 消息类型 TextMessage | BinaryMessage
        ReadBufferSize  int                      // 读取缓存大小
        WriteBufferSize int                      // 写入缓存大小
        CheckOrigin     func(*http.Request) bool // 检查跨域来源是否允许建立连接        
        Logger          loggerI                  // 用于打印内部产生日志
}
    Config 配置项

type Conn struct {
        Socket *websocket.Conn // 连接

        // Has unexported fields.
}
    Conn 连接实例

func (c *Conn) Destroy() error
    Destroy 销毁该连接

func (c *Conn) Init()
    Init 初始化该连接

func (c *Conn) Send(msg *Message) error
    Send 往该连接发送数据

type Message struct {
        SID     uint64
        Payload []byte

        Socket  *Conn
        MsgType int
        // Has unexported fields.
}
    Message ws 接收到的消息

func (m *Message) Response(v []byte) error
    Response 在发送本消息的当前连接发送数据

type WS struct {
        Clients  map[uint64]*Conn
        Upgrader *websocket.Upgrader

        MsgType int

        // Has unexported fields.
}
    WS ws 连接

func New(conf *Config) *WS
    New 新建 websocket 服务

func (ws *WS) HTTPHandler(w http.ResponseWriter, r *http.Request)
    HTTPHandler 给 http 控制器绑定使用

func (ws *WS) Receive() <-chan *Message
    Receive 获取接收数据的 chan

func (ws *WS) Send(msg *Message) error
    Send 发送数据

```