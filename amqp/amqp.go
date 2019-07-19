package amqp

import (
	"errors"
	"github.com/streadway/amqp"
)

type Config struct {
	Addr         string
	ExchangeName string
	Exchange     *Exchange
	Consumer     *Consumer
}

// Init 初始化
func New(conf *Config) (*MQ, error) {
	if conf.Exchange == nil && conf.ExchangeName == "" {
		return nil, errors.New("exchange must defined")
	}

	if conf.Exchange == nil && conf.ExchangeName != "" {
		conf.Exchange = defaultExchange(conf.ExchangeName)
	}

	if conf.Consumer == nil {
		conf.Consumer = defaultConsumer()
	}

	conn, err := amqp.Dial(conf.Addr)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	m := &MQ{conn, channel, conf.Exchange, conf.Consumer}

	err = m.Init()
	if err != nil {
		return nil, err
	}

	return m, nil
}
