package amqp

import "github.com/streadway/amqp"

func defaultConsumer() *Consumer {
	return &Consumer{"", true, false, false, false, nil}
}

// Consumer 定义消费者选项
type Consumer struct {
	Name      string
	AutoAck   bool // 自动确认
	Exclusive bool //
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}
