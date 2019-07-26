package tcp

import (
	"testing"
)

func TestServer(t *testing.T) {
	server, err := NewServer(&Config{
		Network: "tcp",
		Addr:    ":6600",
	})

	server.UseReceive(func(bt []byte) []byte {

	})

	server.UseSend(func(data interface{}) []byte {

	})
	if err != nil {
		panic(err)
	}

	server.Connect()
}
