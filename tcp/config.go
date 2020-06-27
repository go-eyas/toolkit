package tcp

// Config 配置项
type Config struct {
  Addr    string                                // tcp 地址，在客户端使用为需要连接的地址，在服务端使用为监听的地址
  Network string                                // tcp 的网络类型，可选值为 "tcp", "tcp4", "tcp6", "unix" or "unixpacket"
  Packer  func([]byte) ([]byte, error)          // tcp 数据包的封装函数，传入的数据是需要发送的业务数据，返回发送给 tcp 的数据
  Parser  func(*Conn, []byte) ([][]byte, error) // 将收到的数据包，根据私有协议转换成业务数据，在这里处理粘包,半包等数据包问题，返回处理好的数据包
  Logger  LoggerI
}
