package tcp

import (
  "errors"
  "net"
)

type connHandler func(*Conn)

// Server 服务器
type Server struct {
  Listener net.Listener     // 服务器监听实例
  Sockets  map[uint64]*Conn // 当前与客户端的连接实例
  Config   *Config          // 配置项

  recChan            chan *Message // 收到的数据，这里是经过 Parser 解析后的数据
  socketCount        uint64        // id 计数器
  createConnHandlers []connHandler // 当有新连接建立时触发函数
  closeConnHandlers  []connHandler // 当有连接关闭时触发函数
}

// NewServer 实例化服务器
func NewServer(conf *Config) (*Server, error) {
  var defaultParsePoll map[uint64][]byte
  if conf.Packer == nil && conf.Parser == nil {
    conf.Packer = Packer
    defaultParsePoll, conf.Parser = Parser()
  } else if conf.Packer != nil || conf.Parser != nil {
    return nil, errors.New("the Packer and Parser must be specified together")
  }

  if conf.Logger == nil {
    conf.Logger = EmptyLogger
  }

  if conf.Network == "" {
    conf.Network = "tcp"
  }

  listener, err := net.Listen("tcp", conf.Addr)
  if err != nil {
    return nil, err
  }

  server := &Server{
    Listener: listener,
    Config:   conf,
    Sockets:  make(map[uint64]*Conn),
    recChan:  make(chan *Message, 2),
    // sendChan: make(chan *Message, 2),
    createConnHandlers: make([]connHandler, 0),
    closeConnHandlers:  make([]connHandler, 0),
  }

  go server.accept()

  // 清理已关闭的连接解析池
  if defaultParsePoll != nil {
    server.HandleClose(func(conn *Conn) {
      if _, ok := defaultParsePoll[conn.ID]; ok {
      }
    })
  }
  return server, nil
}

// 接收新连接
func (sv *Server) accept() {
  for {
    conn, err := sv.Listener.Accept()
    if err != nil {
      continue
    }
    sv.newConn(conn)
  }
}

// Receive 接收数据
func (sv *Server) Receive() <-chan *Message {
  return sv.recChan
}

// Send 发送数据到指定连接实例
func (sv *Server) Send(conn *Conn, data []byte) error {
  return conn.Send(data)
}

// SendConnID 发送数据到指定连接实例ID
func (sv *Server) SendConnID(id uint64, data []byte) error {
  conn, ok := sv.Sockets[id]
  if !ok || conn == nil {
    return errors.New("invalid connection")
  }
  return sv.Send(conn, data)
}

// HandleCreate 每当有新连接建立时，触发函数
func (sv *Server) HandleCreate(h connHandler) {
  sv.createConnHandlers = append(sv.createConnHandlers, h)
}

// HandleClose 每当有连接关闭时，触发函数
func (sv *Server) HandleClose(h connHandler) {
  sv.closeConnHandlers = append(sv.closeConnHandlers, h)
}

func (sv *Server) newConn(conn net.Conn) *Conn {
  sv.socketCount++
  c := &Conn{
    ID:     sv.socketCount,
    Conn:   conn,
    server: sv,
  }
  sv.Sockets[c.ID] = c

  // 触发器
  for _, h := range sv.createConnHandlers {
    h(c)
  }

  go c.reader()
  return c
}
