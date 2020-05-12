package amqp

import "sync"

type MQApp struct {
	Client         *MQ
	listenerMu     sync.RWMutex
	listener       map[*Queue][]MQHandler
	listenerRecord map[*Queue]bool
}

func NewApp(config *Config) (*MQApp, error) {
	mq, err := New(config)
	return &MQApp{
		Client:         mq,
		listener:       make(map[*Queue][]MQHandler),
		listenerRecord: make(map[*Queue]bool),
	}, err
}

// 监听队列触发函数
func (mq *MQApp) On(queue *Queue, handler ...MQHandler) {
	mq.listenerMu.RLock()
	handlers, ok := mq.listener[queue]
	mq.listenerMu.RUnlock()

	if !ok {
		handlers = handler
	} else {
		handlers = append(handlers, handler...)
	}

	mq.listenerMu.Lock()
	mq.listener[queue] = handlers
	mq.listenerMu.Unlock()

	mq.startListen(queue)
}

func (mq *MQApp) Route(routes map[*Queue]MQHandler) {
	for q, handler := range routes {
		mq.On(q, handler)
	}
}

func (mq *MQApp) startListen(queue *Queue) {
	_, ok := mq.listenerRecord[queue]

	// 之前已经开始了
	if ok {
		return
	}

	// 开始监听
	go func(queue *Queue) {
		ch, err := mq.Client.Sub(queue)
		if err != nil {
			return
		}
		mq.listenerMu.Lock()
		mq.listenerRecord[queue] = true
		mq.listenerMu.Unlock()
		for msg := range ch {

			ctx := &MQContext{
				Request: msg,
				Client:  mq.Client,
				App:     mq,
			}
			handlers := mq.listener[queue]
			go func() {
				// TODO defer error
				for _, h := range handlers {
					h(ctx)
				}
			}()
		}
	}(queue)
}

func (mq *MQApp) Pub(q *Queue, msg *Message) error {
	return mq.Client.Pub(q, msg)
}
