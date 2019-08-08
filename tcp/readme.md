# TCP 

基于底层的 tcp 封装，有 客户端和服务端

# TIPS： 封装未完成，请勿使用

# 客户端 Client


```go
import "github.com/go-eyas/toolkit/tcp"

client, err := tcp.NewClient(&tcp.Config{
	Network: "tcp",
	Addr:    ":6600",

	// 解析私有协议为结构体，如果当前没有解析到，返回nil
	Parser: func(conn *tcp.Conn, bt []byte) (interface{}, error) {
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
err = client.Send(&tcp.Message{
	Data: []byte("hello world1"),
})
```

# 服务端 Server

```go
import "github.com/go-eyas/toolkit/tcp"

server, err := tcp.NewServer(&tcp.Config{
	Network: "tcp",
	Addr:    ":6600",

	// 解析私有协议为结构体，如果当前没有解析到，返回nil
	Parser: func(conn *tcp.Conn, bt []byte) (interface{}, error) {
		return nil, nil
	},
	
	// 将数据转换成字节数组，发送时就发该段数据
	Packer: func(data interface{}) ([]byte, error) {
		return nil, nil
	},
})

if err != nil {
	panic(err)
}

ch := server.Receive()

for data := range ch {
	fmt.Printf("receive: %v", data.Data)
	err := data.Response(map[string]interface{}{
		"hello": "world",
	})
	if err != nil {
		panic(err)
	}
}
```