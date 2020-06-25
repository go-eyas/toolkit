package tcp

import (
	"errors"
	"net"
	"time"
)

type Client struct {
	Conn     *Conn
	Config   *Config
	recChan  chan *Message // 收到的数据，这里是经过 Parser 解析后的数据
	socketCount uint64
	autoReconnect bool

	createConnHandlers []connHandler // 当有新连接建立时触发函数
	closeConnHandlers []connHandler // 当有连接关闭时触发函数

	closeNotify chan *Conn // 连接关闭时通知通道
}

func NewClient(conf *Config) (*Client, error) {
	var defaultParsePoll map[uint64][]byte
	if conf.Packer == nil && conf.Parser == nil {
		conf.Packer = Packer
		defaultParsePoll, conf.Parser = Parser()
	} else if conf.Packer != nil || conf.Parser != nil {
		return nil ,errors.New("the Packer and Parser must be specified together")
	}


	client := &Client{
		autoReconnect: true,
		Config:   conf,
		recChan:  make(chan *Message, 2),
		closeNotify: make(chan *Conn, 0),
	}

	// 连接关闭了通知一下
	client.HandleClose(func(conn *Conn) {
		delete(defaultParsePoll, conn.ID)
		if client.autoReconnect {
			client.closeNotify <- conn
		}
	})

	err := client.connect()

	if err != nil {
		return nil, err
	}

	go client.reconnect()

	return client, nil
}

func (c *Client) connect() error {
	dial, err := net.Dial(c.Config.Network, c.Config.Addr)
	if err != nil {
		return err
	}
	c.socketCount++
	conn := &Conn{Conn: dial, ID: c.socketCount}
	c.Conn = conn
	conn.client = c

	go c.reader()
	return nil
}

func (c *Client) reconnect() {
	if !c.autoReconnect {
		close(c.closeNotify)
		return
	}
	<- c.closeNotify
	// conn := <- c.closeNotify
	// fmt.Printf("conn %d is close, retrying...\n", conn.ID)
	for {
		time.Sleep(1 * time.Second)
		err := c.connect()
		if err != nil {
			// fmt.Printf("reconnect fail: %v\n", err)
		} else {
			// fmt.Printf("reconnect ok: \n")
			go c.reconnect()
			break
		}
	}

}

func (c *Client) reader() {
	for _, h := range c.createConnHandlers {
		h(c.Conn)
	}
	c.Conn.reader()
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

func (c *Client) Destroy() error {
	c.autoReconnect = false
	return c.Conn.Destroy()
}