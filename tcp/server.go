package tcp

import (
	"net"
)

// Server 服务器
type Server struct {
	Listener net.Listener
	Sockets  map[uint64]*Conn

	Config *Config

	// Packer func(interface{}) ([]byte, error)        // 将传入的对象数据，根据私有协议封装成字节数组，用于发送到tcp连接
	// Parser func(*Conn, []byte) (interface{}, error) // 将收到的数据包，根据私有协议转换成具体数据，在这里处理粘包,半包等数据包问题，返回自定义的数据

	recChan  chan *Message // 收到的数据，这里是经过 Parser 解析后的数据
	sendChan chan *Message // 待发送的数据，这里原始数据，下一步会将里面的数据给 Packer 处理好，然后发送出去

	socketCount uint64 // id 计数器
}

func NewServer(conf *Config) (*Server, error) {
	if conf.Packer == nil {
		conf.Packer = Packer
	}
	if conf.Parser == nil {
		conf.Parser = Parser
	}
	listener, err := net.Listen(conf.Network, conf.Addr)
	if err != nil {
		return nil, err
	}

	server := &Server{
		Listener: listener,
		Sockets:  make(map[uint64]*Conn),
		Config:   conf,
		recChan:  make(chan *Message, 2),
		sendChan: make(chan *Message, 2),
	}

	go server.Accept()

	return server, nil
}

func (sv *Server) Accept() {
	for {
		conn, err := sv.Listener.Accept()
		if err == nil {
			sv.newConn(conn)
		}
	}
}

func (sv *Server) Receive() <-chan *Message {
	return sv.recChan
}

func (sv *Server) Send(conn *Conn, msg *Message) error {
	return conn.Send(msg)
}

func (sv *Server) newConn(conn net.Conn) {
	sv.socketCount++
	c := &Conn{
		ID:     sv.socketCount,
		Conn:   conn,
		server: sv,
	}
	sv.Sockets[sv.socketCount] = c
	go c.reader()
}
