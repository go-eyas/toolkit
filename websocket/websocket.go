package websocket

import (
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
	Logger          loggerI                  // 用于打印内部产生日志
}

// New 新建 websocket 服务
func New(conf *Config) *WS {
	if conf.MsgType == 0 {
		conf.MsgType = websocket.TextMessage
	}
	if conf.CheckOrigin == nil {
		conf.CheckOrigin = func(r *http.Request) bool { return true }
	}

	if conf.Logger != nil {
		logger = conf.Logger
	}

	ws := &WS{
		MsgType: conf.MsgType,
		Clients: make(map[uint64]*Conn),
		recC:    make(chan *Message, 1024),
		sendC:   make(chan *Message, 1024),
	}

	ws.Upgrader = &websocket.Upgrader{
		ReadBufferSize:  conf.ReadBufferSize,
		WriteBufferSize: conf.WriteBufferSize,
		CheckOrigin:     conf.CheckOrigin,
	}
	logger.Info("websocket: init websocket")

	return ws
}

// WS ws 连接
type WS struct {
	Clients  map[uint64]*Conn
	Upgrader *websocket.Upgrader
	id       uint64
	MsgType  int

	recC  chan *Message
	sendC chan *Message
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
		// RecChannel:  make(chan *Message, 2),
		// SendChannel: make(chan *Message, 2),
		ws: ws,
		// msgType:     ws.MsgType,
		id: ws.id,
	}

	logger.Infof("websocket: new websocket connect create: %d", conn.id)

	connMu.Lock()
	ws.Clients[conn.id] = conn
	connMu.Unlock()

	conn.Init()

	defer ws.destroyConn(ws.id)
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

// destroyConn 销毁连接
func (ws *WS) destroyConn(cid uint64) {
	conn, ok := ws.Clients[ws.id]
	if !ok {
		return
	}
	conn.Destroy()

	connMu.Lock()
	delete(ws.Clients, cid)
	connMu.Unlock()
	logger.Info("websocket: destroy ws connect")
}
