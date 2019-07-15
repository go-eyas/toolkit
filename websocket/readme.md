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

[API 文档](https://gowalker.org/github.com/go-eyas/toolkit/websocket)