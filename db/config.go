package db

// Config 数据库配置项
type Config struct {
	Driver string `yaml:"driver" json:"driver" toml:"driver" env:"DB_DRIVER"`
	URI    string `yaml:"uri" json:"uri" toml:"uri" env:"DB_URI"`
	Debug  bool
	Logger Logger
}

// Logger 日志对象
type Logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}
