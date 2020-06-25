package tcp

import (
	"errors"
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

func (conn *Conn) IsServer() bool {
	return conn.server != nil
}

func (conn *Conn) IsClient() bool {
	return conn.client != nil
}

func (conn *Conn) reader() {
	var parser func(*Conn, []byte) ([]interface{}, error)
	if conn.IsClient() {
		parser = conn.client.Config.Parser
	} else if conn.IsServer() {
		parser = conn.server.Config.Parser
	}
	for {
		_buf := make([]byte, 1024)
		buflen, err := conn.Conn.Read(_buf)
		if err != nil {
			// 数据异常，马上断开连接
			conn.Destroy()
			break
		}
		buf := _buf[:buflen]
		body, err := parser(conn, buf)
		if err != nil {
			// 解析异常，断开连接
			conn.Destroy()
			break
		}
		for _, body := range body {
			msg := &Message{
				Data: body,
				Conn: conn,
			}
			if conn.IsServer() {
				conn.server.recChan <- msg
			} else if conn.IsClient() {
				conn.client.recChan <- msg
			}
		}
	}
}

func (conn *Conn) Send(msg *Message) error {
	conn.writeMu.Lock()
	defer conn.writeMu.Unlock()
	var err error
	var pack []byte

	if conn.IsClient() {
		pack, err = conn.client.Config.Packer(msg.Data)
	} else if conn.IsServer() {
		pack, err = conn.server.Config.Packer(msg.Data)
	} else {
		err = errors.New("the connection is invalid")
	}

	if err != nil {
		return err
	}
	_, err = conn.Conn.Write(pack)
	return err
}

func (conn *Conn) Destroy() error {
	if conn.IsServer() {
		for _, h := range conn.server.closeConnHandlers {
			h(conn)
		}

		delete(conn.server.Sockets, conn.ID)
	} else if conn.IsClient() {
		for _, h := range conn.client.closeConnHandlers {
			h(conn)
		}
	}
	return conn.Conn.Close()
}
