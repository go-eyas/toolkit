package amqp

import (
	"errors"
)

const (
	ExchangeDirect  = "direct"  // 直连交换机
	ExchangeFanout  = "fanout"  // 扇形交换机
	ExchangeTopic   = "topic"   // 主题交换机
	ExchangeHeaders = "headers" // 头交换机
)

// Config 配置项
// ExchangeName 和 Exchange 二选一，用于指定发布和订阅时使用的交换机
type Config struct {
	Addr         string    // rabbitmq 地址
	ExchangeName string    // 使用该值创建一个直连的交换机
	Exchange     *Exchange // 自定义默认交换机
	Consumer     *Consumer // 在定于队列时，作为消费者使用的参数
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

	m := &MQ{Addr: conf.Addr, Exchange: conf.Exchange, Consumer: conf.Consumer}

	err := m.Init()
	if err != nil {
		return nil, err
	}

	return m, nil
}
