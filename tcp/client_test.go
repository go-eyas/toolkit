package tcp

import (
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client, err := NewClient(&Config{
		Network: "tcp",
		Addr:    ":6600",

		// 解析私有协议为结构体，如果当前没有解析到，返回nil
		Parser: func(conn *Conn, bt []byte) (interface{}, error) {
			return nil, nil
		},

		// 将数据转换成字节数组，发送时就发该段数据
		Packer: func(data interface{}) ([]byte, error) {
			return nil, nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// 发送数据
	err = client.Send(&Message{
		Data: []byte("hello world1"),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = client.Send(&Message{
		Data: []byte("hello world2"),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = client.Send(&Message{
		Data: []byte("hello world3"),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = client.Send(&Message{
		Data: []byte("hello world4"),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = client.Send(&Message{
		Data: []byte("hello world5"),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = client.Send(&Message{
		Data: []byte("hello world"),
	})
	if err != nil {
		t.Fatal(err)
	}
	<-time.After(3 * time.Second)

	// ch, err := client.Receive()
	//
	// for data := range ch {
	// 	fmt.Println("receive data:", data.Body)
	// 	data.Response(map[string]interface{}{
	// 		"msg": "I'm fine.",
	// 	})
	// }

}
