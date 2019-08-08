package tcp

import (
	"fmt"
	"net"
	"sync"
)

type Conn struct {
	writeMu sync.Mutex
	ID      uint64
	Conn    net.Conn
	server  *Server
	client  *Client
}

func (conn *Conn) reader() {
	for {
		_buf := make([]byte, 1024)
		buflen, err := conn.Conn.Read(_buf)
		fmt.Println("rec buf", buflen)
		if err != nil {
			// 数据异常，马上断开连接
			fmt.Println(err)
			conn.Destroy()
			break
		}
		buf := _buf[:buflen]
		body, err := conn.server.Config.Parser(conn, buf)
		fmt.Printf("rec: %s\n", string(body.([]byte)))
		if err != nil {
			data := &Message{
				Data: body,
				Conn: conn,
			}
			conn.server.recChan <- data
		}
	}
}

func (conn *Conn) Send(msg *Message) error {
	conn.writeMu.Lock()
	defer conn.writeMu.Unlock()
	var bt []byte
	var err error
	if conn.server != nil {
		bt, err = conn.server.Config.Packer(msg.Data)
	} else if conn.client != nil {
		bt, err = conn.client.Config.Packer(msg.Data)
	}
	if err != nil {
		return err
	}
	_, err = conn.Conn.Write(bt)
	return err
}

func (conn *Conn) Destroy() {
	if conn.server != nil {
		delete(conn.server.Sockets, conn.ID)
	}
	_ = conn.Conn.Close()
}
