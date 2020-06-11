package db

import (
	"fmt"
	"github.com/go-eyas/toolkit/log"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
	"github.com/novalagung/gubrak"

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
		msgs := gorm.LogFormatter(v...)
		l.logger.Debug("SQL [", v[5], " rows][", tm.String(), "]: ", msgs[3])
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

	var logger *gormLogger
	if conf.Logger != nil {
		logger = &gormLogger{conf.Logger}
	} else {
		logger = &gormLogger{log.SugaredLogger}
	}
	db.SetLogger(logger)
	if conf.Debug {
		db.LogMode(conf.Debug)
	}

	return db, nil
}

func GormViewMigrate(db *gorm.DB, v ...ViewModel) {
	for _, m := range v {
		if db.HasTable(m) {
			continue
		}
		var tags []string
		scope := db.NewScope(m)
		ms := scope.GetModelStruct()
		for _, field := range ms.StructFields {
			if !field.IsIgnored {
				tags = append(tags, scope.Quote(field.DBName))
			}
		}

		// 去掉重复字段
		_tags, _ := gubrak.Filter(tags, func(s string) bool {
			res, _ := gubrak.Find(tags, func(i string) bool {
				index := strings.LastIndex(i, "."+s)
				return index != -1
			})
			return res == nil
		})
		tags = _tags.([]string)

		db.Exec(fmt.Sprintf("CREATE VIEW %v AS SELECT %s %s", scope.QuotedTableName(), strings.Join(tags, ","), m.From()))
	}
}
