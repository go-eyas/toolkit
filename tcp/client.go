package tcp

import (
	"errors"
	"net"
)

type Client struct {
	Conn     *Conn
	Config   *Config
	recChan  chan *Message // 收到的数据，这里是经过 Parser 解析后的数据
	// sendChan chan *Message // 待发送的数据，这里原始数据，下一步会将里面的数据给 Packer 处理好，然后发送出去

	createConnHandlers []connHandler // 当有新连接建立时触发函数
	closeConnHandlers []connHandler // 当有连接关闭时触发函数
}

func NewClient(conf *Config) (*Client, error) {
	if conf.Packer == nil && conf.Parser == nil {
		conf.Packer = Packer
		_, conf.Parser = Parser()
	} else if conf.Packer != nil || conf.Parser != nil {
		return nil, errors.New("the Packer and Parser must be specified together")
	}

	dial, err := net.Dial(conf.Network, conf.Addr)
	if err != nil {
		return nil, err
	}

	conn := &Conn{Conn: dial}
	client := &Client{
		Conn:     conn,
		Config:   conf,
		recChan:  make(chan *Message, 2),
		// sendChan: make(chan *Message, 2),
	}
	conn.client = client

	go client.Conn.reader()

	return client, nil
}

// TODO
func (c *Client) reconnect() {

}

func (c *Client) reader() {
	go c.Conn.reader()
}

func (c *Client) HandleCreate(h connHandler) {
	c.createConnHandlers = append(c.createConnHandlers, h)
}

func (c *Client) HandleClose(h connHandler) {
	c.closeConnHandlers = append(c.closeConnHandlers, h)
}

func (c *Client) Receive() <-chan *Message {
	return c.recChan
}

func (c *Client) Send(msg *Message) error {
	return c.Conn.Send(msg)
}
