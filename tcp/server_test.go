package tcp

import (
	"testing"
)

func TestServer(t *testing.T) {
	server, err := NewServer(&Config{
		Network: "tcp",
		Addr:    ":6600",

		// // 解析私有协议为结构体，如果当前没有解析到，返回 error
		// Parser: func(conn *Conn, bt []byte) (interface{}, error) {
		// 	return nil, nil
		// },
		//
		// // 将数据转换成字节数组，发送时就发该段数据，如果解析错误返回 error
		// Packer: func(data interface{}) ([]byte, error) {
		// 	return nil, nil
		// },
	})

	if err != nil {
		panic(err)
	}

	ch := server.Receive()

	for data := range ch {
		t.Logf("receive: %v", string(data.Data))
		err := data.Response([]byte(`server receive: ` + string(data.Data)))
		if err != nil {
			t.Log(err)
		}
	}

}
