package amqp

import (
	"encoding/json"
)

type Message struct {
	ContentType string
	mq          *MQ
	Queue       *Queue
	Data        []byte
}

func (m *Message) contentType() string {
	if m.ContentType == "" {
		return "text/plain"
	}
	return m.ContentType
}

func (m *Message) JSON(v interface{}) error {
	return json.Unmarshal(m.Data, v)
}

func (m *Message) ReplyTo(msg *Message) error {
	return m.mq.Pub(m.Queue.ReplyTo, msg)
}
