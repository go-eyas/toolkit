package amqp

import (
	"encoding/json"
)

// Message 消息体
type Message struct {
	ContentType string // 消息类型
	Queue       *Queue // 来自于哪个队列
	Data        []byte // 消息数据
	mq          *MQ
}

func (m *Message) contentType() string {
	if m.ContentType == "" {
		return "text/plain"
	}
	return m.ContentType
}

// JSON 以 json 解析消息体的数据为指定结构体
func (m *Message) JSON(v interface{}) error {
	return json.Unmarshal(m.Data, v)
}

// ReplyTo 给回复的队列发送消息
func (m *Message) ReplyTo(msg *Message) error {
	return m.mq.Pub(m.Queue.ReplyTo, msg)
}
