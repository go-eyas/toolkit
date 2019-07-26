package tcp

import "net"

type ServerConn struct {
	ID       uint
	Conn     net.Conn
	recChan  chan []byte
	sendChan chan []byte
}
