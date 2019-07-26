# 配置文件

 * 支持 YAML, JSON, TOML, 环境变量 设置配置项的值
 * 支持多文件覆盖配置
 * 支持默认值
 * 修改文件后自动重载

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

## 配置文件加载顺序

后面的会覆盖前面的，假设 file 参数传的是 `config`， shell 的环境变量 `CONFIG_ENV=dev`

```
../config.dev.yml
../config.yml
config.dev.yml
config.yml
../config.dev.json
../config.json
config.dev.json
config.json
../config.dev.toml
../config.toml
config.dev.toml
config.toml
环境变量
```

> 如果环境变量 `CONFIG_ENV` 没有设置，则默认为 `local`

# godoc

[API 文档](https://gowalker.org/github.com/go-eyas/toolkit/config)