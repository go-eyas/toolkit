package amqp

import (
	"fmt"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type MQ struct {
	Addr        string
	Client      *amqp.Connection
	Channel     *amqp.Channel
	Exchange    *Exchange
	Consumer    *Consumer
	notifyClose chan *amqp.Error
	subQueues   []*Queue // 已注册为消费者的通道

}

// Init 初始化
// 1. 初始化交换机
func (mq *MQ) Init() error {
	mq.subQueues = []*Queue{}

	err := mq.connect()
	if err != nil {
		return err
	}

	// 初始化默认交换机
	err = mq.ExchangeDeclare(mq.Exchange)
	if err != nil {
		return err
	}

	return nil
}

func (mq *MQ) connect() error {
	conn, err := amqp.Dial(mq.Addr)
	if err != nil {
		return err
	}
	mq.Client = conn
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	mq.Channel = channel

	// 重连后重新注册消费者
	for _, q := range mq.subQueues {
		q.IsDeclare = false
		q.exchange = nil
		q.q = nil
		mq.bindMQChan(q)
	}

	// 断线重连
	go mq.reconnect()

	return nil
}

func (mq *MQ) reconnect() {
	mq.notifyClose = make(chan *amqp.Error)
	mq.Channel.NotifyClose(mq.notifyClose)

	for {
		select {
		case <-mq.notifyClose:
			fmt.Println("rabbitmq connection is close, retrying...")
			mq.Client.Close()
			<-time.After(500 * time.Millisecond) // 隔 500ms 重连一次
			mq.connect()
			break
		}
	}
}

var subMu sync.Mutex

// Sub 定于队列消息
// q 队列
//return 接收消息的通道 ， 错误对象
func (mq *MQ) Sub(q *Queue) (<-chan *Message, error) {
	subMu.Lock()
	defer subMu.Unlock()

	mq.subQueues = append(mq.subQueues, q)

	// 初始化接收通道
	if q.consumerChan == nil {
		q.consumerChan = make(chan *Message, 2)
	}

	mq.bindMQChan(q)

	return q.consumerChan, nil
}

var bindMu sync.Mutex

// 将 mq 通道绑到队列通道中
func (mq *MQ) bindMQChan(q *Queue) error {
	bindMu.Lock()
	defer bindMu.Unlock()
	// 定义队列
	if !q.IsDeclare {
		err := mq.QueueDeclare(q)
		if err != nil {
			return err
		}
	}

	e := mq.Exchange

	// 绑定交换机
	if q.exchange != e {
		err := mq.QueueBind(q, e)
		if err != nil {
			return err
		}
	}
	msgChan, err := mq.Channel.Consume(
		q.Name,
		mq.Consumer.Name,
		mq.Consumer.AutoAck,
		mq.Consumer.Exclusive,
		mq.Consumer.NoLocal,
		mq.Consumer.NoWait,
		mq.Consumer.Args,
	)

	if err != nil {
		return err
	}

	go func(ch chan<- *Message) {
		for d := range msgChan {
			msg := &Message{
				ContentType: d.ContentType,
				mq:          mq,
				Queue:       q,
				Data:        d.Body,
			}
			ch <- msg
		}
	}(q.consumerChan)

	return nil
}

var pubMu sync.Mutex

// Pub 给队列发送消息,
// q 队列,
// msg 消息,
// exchanges 交换机，可以用多个交换机多次发送，默认使用初始化时指定的交换机
func (mq *MQ) Pub(q *Queue, msg *Message, exchanges ...*Exchange) error {
	pubMu.Lock()
	defer pubMu.Unlock()

	// 定义队列
	if !q.IsDeclare {
		err := mq.QueueDeclare(q)
		if err != nil {
			return err
		}
	}

	if len(exchanges) == 0 {
		exchanges = append(exchanges, mq.Exchange)
		// 绑定初始化的交换机
		if q.exchange != mq.Exchange {
			err := mq.QueueBind(q, mq.Exchange)
			if err != nil {
				return err
			}
		}
	} else {
		for _, e := range exchanges {
			if !e.IsDeclare {
				err := mq.ExchangeDeclare(e)
				if err != nil {
					return err
				}
			}
			err := mq.Channel.QueueBind(
				q.Name,
				q.GetKey(),
				e.Name,
				false,
				nil,
			)
			if err != nil {
				return err
			}
		}
	}

	for _, e := range exchanges {
		// 发消息
		err := mq.Channel.Publish(
			e.Name,
			q.GetKey(),
			false,
			false,
			amqp.Publishing{
				ContentType: msg.ContentType,
				ReplyTo:     q.ReplyQueue(),
				Body:        msg.Data,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil

}

func (mq *MQ) QueueDeclare(q *Queue) error {
	queue, err := mq.Channel.QueueDeclare(
		q.Name,
		q.Durable,
		q.AutoDelete,
		q.Exclusive,
		q.NoWait,
		q.Args,
	)
	if err != nil {
		return err
	}
	q.q = &queue
	q.IsDeclare = true
	return nil
}

func (mq *MQ) ExchangeDeclare(e *Exchange) error {
	if e.IsDeclare {
		return nil
	}
	err := mq.Channel.ExchangeDeclare(
		e.Name,
		e.Kind,
		e.Durable,
		e.AutoDelete,
		e.Internal,
		e.NoWait,
		e.Args,
	)
	if err == nil {
		e.IsDeclare = true
	}
	return err
}

func (mq *MQ) QueueBind(q *Queue, e *Exchange) error {
	if !e.IsDeclare {
		mq.ExchangeDeclare(e)
	}
	err := mq.Channel.QueueBind(
		q.Name,
		q.GetKey(),
		e.Name,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	q.exchange = e
	return nil
}
