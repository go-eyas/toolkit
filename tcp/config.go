package tcp

// Config 配置项
type Config struct {
	Addr    string
	Network string
	Packer  func(interface{}) ([]byte, error)        // 将传入的对象数据，根据私有协议封装成字节数组，用于发送到tcp连接
	Parser  func(*Conn, []byte) ([]interface{}, error) // 将收到的数据包，根据私有协议转换成具体数据，在这里处理粘包,半包等数据包问题，返回自定义的数据
}
