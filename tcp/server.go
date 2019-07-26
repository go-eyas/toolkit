package tcp

import (
	"net"
)

// Server 服务器
type Server struct {
	listener net.Listener
	sockets  map[int]net.Conn

	pkgParser PackageParser
	recChan   chan interface{}
	sendChan  chan []byte

	socketCount int
}

func NewServer(conf *Config) (*Server, error) {
	listener, err := net.Listen(conf.Network, conf.Addr)
	if err != nil {
		return nil, err
	}
	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}
	server := &Server{
		listener:  listener,
		conn:      conn,
		pkgParser: conf.PkgParser,
		recChan:   make(chan interface{}, 2),
		sendChan:  make(chan []byte, 2),
	}

	return server, nil
}

func (sv *Server) reader() {
	for {
		_buf := make([]byte, 1024)
		buflen, err := sv.conn.Read(_buf)
		buf := _buf[:buflen]
		data := sv.pkgParser.Parser(buf)
		sv.recChan <- data
	}
}

func (sv *Server) writer() {
	for data := range sv.sendChan {
		buf := sv.pkgParser.Packer(data)
		res, err := sv.conn.Write(buf)
	}
}
