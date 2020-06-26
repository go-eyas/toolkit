package tcp

import (
	"fmt"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client, err := NewClient(&Config{
		Network: "tcp",
		Addr:    ":6600",

		// 解析私有协议为结构体，如果当前没有解析到，返回nil
		// Parser: func(conn *Conn, bt []byte) (interface{}, error) {
		// 	return nil, nil
		// },
		//
		// // 将数据转换成字节数组，发送时就发该段数据
		// Packer: func(data interface{}) ([]byte, error) {
		// 	return nil, nil
		// },
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("connect ok")

	// 发送数据
	t.Log("send 1")
	err = client.Send([]byte("hello world1"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("send 2")
	err = client.Send([]byte("hello world2"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("send 3")
	err = client.Send([]byte("hello world3"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("send 4")
	err = client.Send([]byte("hello world4"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("send 5")
	err = client.Send([]byte("hello world5"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("send 6")
	err = client.Send([]byte("hello world6"))
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for i := 0; true; i++ {
			time.Sleep(2 * time.Second)
			err = client.Send([]byte(fmt.Sprintf("auto send %d", i)))

			if i > 20 {
				client.Destroy()
			}
			if err != nil {
				t.Log(err)
			}
		}
	}()


	ch := client.Receive()

	for data := range ch {
		fmt.Println("receive data:", string(data.Data))
		// data.Response(map[string]interface{}{
		// 	"msg": "I'm fine.",
		// })
	}

}
