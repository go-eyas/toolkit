# 配置文件

 * 支持 YAML, JSON, TOML, 环境变量 设置配置项的值
 * 支持自动重载
 * 支持多文件覆盖配置
 * 支持默认值

# 使用

```go
import "github.com/go-eyas/toolkit/config"

type Config struct {
	IsParseJSON      bool `json:"isParseJson"`
	IsParseToml      bool `toml:"isParseToml"`
	IsParseYml       bool `yaml:"isParseYml"`
	IsParseLocalToml bool `toml:"isParseLocalToml"`

	Ext string 
	Env string `default:"shell" env:"APP_ENV" required:"true"`
	Obj struct {
		Array   []int
		Boolean bool
	}
}

func main() {
  conf := &Config{}
  err := config.Init("test/config", conf)
}

```