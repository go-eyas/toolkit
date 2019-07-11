package config

import (
	"testing"
)

type ConfigT struct {
	IsParseJSON      bool `json:"isParseJson"`
	IsParseToml      bool `toml:"isParseToml"`
	IsParseYml       bool `yaml:"isParseYml"`
	IsParseLocalToml bool `toml:"isParseLocalToml"`

	Ext string
	Env string
	Obj struct {
		Array   []int
		Boolean bool
	}
}

func TestConfig(t *testing.T) {
	files := parseFiles("test/config")
	t.Logf("valid files: %+v", files)

	conf := &ConfigT{}
	err := Init("test/config", conf)

	if err != nil {
		panic(err)
	}
	t.Logf("config: %+v", conf)
	if conf.IsParseJSON && conf.IsParseToml && conf.IsParseYml && conf.IsParseLocalToml {
		t.Log("parse config success")
	} else {
		panic("parse config error")
	}
}
