package websocket

import (
  "errors"
)

// Message ws 接收到的消息
type Message struct {
  SID      uint64
  Payload  []byte
  ws       *WS
  Socket   *Conn
  MsgType  int
}

func (m Message) clone() *Message {
  m1 := m
  return &m1
}

func (m *Message) writer() error {
  if m.Socket.isClose {
    return errors.New("socket is close")
  }
  err := m.Socket.Socket.WriteMessage(m.MsgType, m.Payload)
  if err != nil {
    return err
  }
  return nil
}

// Response 在发送本消息的当前连接发送数据
func (m *Message) Response(v []byte) error {
  resMsg := m.clone()
  resMsg.Payload = v
  return resMsg.writer()
}
