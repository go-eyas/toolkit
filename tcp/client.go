package tcp

import (
	"net"
)

type Client struct {
	Conn     *Conn
	Config   *Config
	recChan  chan *Message // 收到的数据，这里是经过 Parser 解析后的数据
	sendChan chan *Message // 待发送的数据，这里原始数据，下一步会将里面的数据给 Packer 处理好，然后发送出去
}

func NewClient(conf *Config) (*Client, error) {
	if conf.Packer == nil {
		conf.Packer = Packer
	}
	if conf.Parser == nil {
		conf.Parser = Parser
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
		sendChan: make(chan *Message, 2),
	}
	conn.client = client

	return client, nil
}

func (c *Client) Receive() <-chan *Message {
	return c.recChan
}

func (c *Client) Send(msg *Message) error {
	return c.Conn.Send(msg)
}
