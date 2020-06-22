package websocket

import (
	"github.com/gorilla/websocket"
)

// Conn 连接实例
type Conn struct {
	Socket  *websocket.Conn // 连接
	ws      *WS             // ws 服务
	isClose bool
	ID      uint64
}

// Init 初始化该连接
func (c *Conn) Init() {
	c.reader()
}

func (c *Conn) reader() error {
	for {
		mType, mRaw, err := c.Socket.ReadMessage()
		if err != nil {
			return err
		}
		logger.Infof("websocket: receive data=%s", string(mRaw))
		msg := &Message{c.ID, mRaw, c.ws, c, mType}
		c.ws.recC <- msg
	}
}

// Send 往该连接发送数据
func (c *Conn) Send(msg *Message) error {
	m := &(*msg)
	m.ws = c.ws
	m.Socket = c
	return m.writer()
}

// Destroy 销毁该连接
func (c *Conn) Destroy() error {
	if c.isClose {
		return nil
	}
	err := c.Socket.Close()
	if err != nil {
		return err
	}
	c.isClose = true

	return nil
}
