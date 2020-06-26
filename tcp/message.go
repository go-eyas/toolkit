package tcp

type Message struct {
	Data []byte
	Conn *Conn
}

func (m *Message) Response(data []byte) error {
	return m.Conn.Send(data)
}
