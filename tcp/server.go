package tcp

import (
	"errors"
	"net"
)

type connHandler func(*Conn)

// Server 服务器
type Server struct {
	Listener net.Listener
	Sockets  map[uint64]*Conn

	Config *Config

	recChan  chan *Message // 收到的数据，这里是经过 Parser 解析后的数据
	// sendChan chan *Message // 待发送的数据，这里原始数据，下一步会将里面的数据给 Packer 处理好，然后发送出去

	socketCount uint64 // id 计数器

	createConnHandlers []connHandler // 当有新连接建立时触发函数
	closeConnHandlers []connHandler // 当有连接关闭时触发函数
}

func NewServer(conf *Config) (*Server, error) {
	var defaultParsePoll map[uint64][]byte
	if conf.Packer == nil && conf.Parser == nil {
		conf.Packer = Packer
		defaultParsePoll, conf.Parser = Parser()
	} else if conf.Packer != nil || conf.Parser != nil {
		return nil ,errors.New("the Packer and Parser must be specified together")
	}

	listener, err := net.Listen(conf.Network, conf.Addr)
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
		closeConnHandlers: make([]connHandler, 0),
	}

	go server.Accept()

	// 清理已关闭的连接解析池
	if defaultParsePoll != nil {
		server.HandleClose(func(conn *Conn) {
			if _, ok := defaultParsePoll[conn.ID]; ok {}
		})
	}
	return server, nil
}

func (sv *Server) Accept() {
	for {
		conn, err := sv.Listener.Accept()
		if err != nil {
			continue
		}
		sv.newConn(conn)
	}
}

func (sv *Server) Receive() <-chan *Message {
	return sv.recChan
}

func (sv *Server) Send(conn *Conn, data interface{}) error {
	msg := &Message{
		Data: data,
		Conn: conn,
	}
	return conn.Send(msg)
}

func (sv *Server) HandleCreate(h connHandler) {
	sv.createConnHandlers = append(sv.createConnHandlers, h)
}

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
