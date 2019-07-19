package amqp

import (
	"github.com/streadway/amqp"
	"sync"
)

type MQ struct {
	Client   *amqp.Connection
	Channel  *amqp.Channel
	Exchange *Exchange
	Consumer *Consumer
}

func (mq *MQ) Init() error {
	err := mq.exchangeDeclare(mq.Exchange)
	if err != nil {
		return err
	}
	return nil
}

type SubHandler func(*Message)

var subMu sync.Mutex

func (mq *MQ) Sub(q *Queue) (<-chan *Message, error) {
	subMu.Lock()
	defer subMu.Unlock()

	// 定义队列
	if !q.IsDeclare {
		err := mq.queueDeclare(q)
		if err != nil {
			return nil, err
		}
	}

	e := mq.Exchange

	// 绑定交换机
	if q.exchange != e {
		err := mq.queueBind(q, e)
		if err != nil {
			return nil, err
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
		return nil, err
	}

	ch := make(chan *Message, 2)

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
	}(ch)

	return ch, nil

}

var pubMu sync.Mutex

func (mq *MQ) Pub(q *Queue, msg *Message) error {
	pubMu.Lock()
	defer pubMu.Unlock()

	// 定义队列
	if !q.IsDeclare {
		err := mq.queueDeclare(q)
		if err != nil {
			return err
		}
	}
	e := mq.Exchange
	// 绑定交换机
	if q.exchange != e {
		err := mq.queueBind(q, e)
		if err != nil {
			return err
		}
	}

	// 发消息
	err := mq.Channel.Publish(
		e.Name,
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: msg.ContentType,
			ReplyTo:     q.replyTo(),
			Body:        msg.Data,
		},
	)

	return err

}

func (mq *MQ) queueDeclare(q *Queue) error {
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

func (mq *MQ) exchangeDeclare(e *Exchange) error {
	err := mq.Channel.ExchangeDeclare(
		e.Name,
		e.Kind,
		e.Durable,
		e.AutoDelete,
		e.Internal,
		e.NoWait,
		e.Args,
	)
	return err
}

func (mq *MQ) queueBind(q *Queue, e *Exchange) error {
	err := mq.Channel.QueueBind(
		q.Name,
		q.Name,
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
