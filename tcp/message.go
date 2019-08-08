package tcp

type Message struct {
	Data interface{}
	Conn *Conn
}

func (m *Message) Response(data interface{}) error {
	msg := &Message{
		Data: data,
		Conn: m.Conn,
	}
	return m.Conn.Send(msg)
}
