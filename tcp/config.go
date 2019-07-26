package tcp

// Config 配置项
type Config struct {
	Addr      string
	Network   string
	PkgParser PackageParser
}