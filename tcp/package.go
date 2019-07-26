package tcp

// PackageParser 封装和解析数据包，实现私有协议
type PackageParser interface {
	Packer(interface{}) []byte // 将传入的对象数据，根据私有协议封装成字节数组，用于发送到tcp连接
	Parser([]byte) interface{} // 将收到的数据包，根据私有协议转换成具体数据，在这里处理粘包,半包等数据包问题，返回自定义的数据
}