package websocket

import (
  "html/template"
  "net/http"
  "sync"

  "github.com/gorilla/websocket"
)

const (
  // 文本消息
  TextMessage = 1

  // 二进制数据消息
  BinaryMessage = 2
)

// Config 配置项
type Config struct {
  MsgType         int                      // 消息类型 TextMessage | BinaryMessage
  ReadBufferSize  int                      // 读取缓存大小
  WriteBufferSize int                      // 写入缓存大小
  CheckOrigin     func(*http.Request) bool // 检查跨域来源是否允许建立连接
  Logger          LoggerI                  // 用于打印内部产生日志
}

// New 新建 websocket 服务
func New(conf *Config) *WS {
  if conf.MsgType == 0 {
    conf.MsgType = websocket.TextMessage
  }
  if conf.CheckOrigin == nil {
    conf.CheckOrigin = func(r *http.Request) bool { return true }
  }

  ws := &WS{
    MsgType:   conf.MsgType,
    Clients:   make(map[uint64]*Conn),
    recC:      make(chan *Message, 1024),
    logger:    conf.Logger,
    createHandlers: make([]EventHandle, 0),
    closeHandlers: make([]EventHandle, 0),
  }

  if ws.logger == nil {
    ws.logger = EmptyLogger
  }

  ws.Upgrader = &websocket.Upgrader{
    ReadBufferSize:  conf.ReadBufferSize,
    WriteBufferSize: conf.WriteBufferSize,
    CheckOrigin:     conf.CheckOrigin,
  }
  ws.logger.Info("websocket: init websocket")

  return ws
}

type EventHandle func(*Conn)

// WS ws 连接
type WS struct {
  Clients  map[uint64]*Conn
  Upgrader *websocket.Upgrader
  id       uint64
  MsgType  int
  recC      chan *Message
  logger LoggerI
  createHandlers []EventHandle
  closeHandlers []EventHandle
}

var connMu sync.RWMutex

// HTTPHandler 给 http 控制器绑定使用
func (ws *WS) HTTPHandler(w http.ResponseWriter, r *http.Request) {
  socket, err := ws.Upgrader.Upgrade(w, r, nil)
  if err != nil {
    return
  }
  ws.id++

  conn := &Conn{
    Socket: socket,
    ws:     ws,
    ID:     ws.id,
  }

  ws.logger.Infof("websocket: new websocket connect create: sid=%d", conn.ID)

  connMu.Lock()
  ws.Clients[conn.ID] = conn
  connMu.Unlock()

  // send init message
  for _, createH := range ws.createHandlers {
    createH(conn)
  }

  conn.Init()

  defer ws.destroyConn(ws.id)
}

var page = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        var addr = document.getElementById("addr").value
        ws = new WebSocket(addr);
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<p>addr: 
<input id="addr" type="text" value="{{.}}">
</p>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))

func (ws *WS) Playground(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("Content-Type", "text/html")
  err := page.Execute(w, "ws://"+r.Host+"/ws")
  if err != nil {
    panic(err)
  }
}

// Receive 获取接收数据的 chan
func (ws *WS) Receive() <-chan *Message {
  return ws.recC
}

// Send 发送数据
func (ws *WS) Send(msg *Message) error {
  m := &(*msg)
  m.ws = ws
  return m.writer()
}

func (ws *WS) HandleClose(fn EventHandle) {
  ws.closeHandlers = append(ws.closeHandlers, fn)
}

func (ws *WS) HandleCreate(fn EventHandle) {
  ws.createHandlers = append(ws.createHandlers, fn)
}

// destroyConn 销毁连接
func (ws *WS) destroyConn(cid uint64) {
  conn, ok := ws.Clients[ws.id]
  if !ok {
    return
  }
  for _, closeH := range ws.closeHandlers {
    closeH(conn)
  }

  conn.Destroy()

  connMu.Lock()
  delete(ws.Clients, cid)
  connMu.Unlock()
  ws.logger.Info("websocket: destroy ws connect")
}
