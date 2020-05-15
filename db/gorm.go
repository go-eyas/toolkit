package db

import (
	"github.com/jinzhu/gorm"
	"time"

	// load drivers
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type gormLogger struct {
	logger Logger
}

func (l *gormLogger) Print(v ...interface{}) {
	var level = v[0]

	if level == "sql" {
		tm := v[2].(time.Duration)
		sql := v[3]
		//不能用log.Println,因为这样log会混乱重合在一起
		l.logger.Debug("SQL [", v[5], " rows][", tm.String(), "]: ", sql, " <-- ", v[4])
	} else {
		l.logger.Debug(v...)
	}
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
