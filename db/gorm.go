package db

import (
	"github.com/jinzhu/gorm"
	// load drivers
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type gormLogger struct {
	logger Logger
}

func (gl *gormLogger) Print(v ...interface{}) {
	gl.logger.Debug(v...)
}

// Gorm 初始化 gorm，返回 gorm 实例
func Gorm(conf *Config) (*gorm.DB, error) {
	db, err := gorm.Open(conf.Driver, conf.URI)
	if err != nil {
		return nil, err
	}

	if conf.Debug {
		db.LogMode(conf.Debug)
		if conf.Logger != nil {
			log := &gormLogger{conf.Logger}
			db.SetLogger(log)
		}
	}

	return db, nil
}
