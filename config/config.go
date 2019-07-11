package config

import (
	"os"

	"github.com/jinzhu/configor"
)

// 自动搜索配置文件 config.xxx 并自动加载，如果配置文件不存在，使用默认配置
// 支持三种配置文件格式
// 并支持 环境变量覆盖配置文件

// 获取有效的
func parseFiles(name string) []string {
	exts := []string{"toml", "json", "yml"}
	env := os.Getenv("CONFIG_ENV")
	if env == "" {
		env = "local"
	}

	filelist := []string{}

	for _, ext := range exts {
		filelist = append(filelist,
			name+"."+ext,
			name+"."+env+"."+ext,
			"../"+name+"."+ext,
			"../"+name+"."+env+"."+ext,
		)
	}

	validFiles := []string{}
	for _, f := range filelist {
		if _, err := os.Stat(f); !os.IsNotExist(err) {
			validFiles = append(validFiles, f)
		}
	}

	return validFiles
}

// Init 初始化配置文件
func Init(file string, v interface{}) error {
	files := parseFiles(file)

	err := configor.Load(v, files...)

	return err
}
