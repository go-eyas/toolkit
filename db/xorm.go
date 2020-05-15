package db

import (
	// load mysql driver
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/log"

	// load postgresql driver
	_ "github.com/lib/pq"

	// load mssql
	_ "github.com/denisenkom/go-mssqldb"
)

type xormLogger struct {
	logger Logger
}

func (xl *xormLogger) Debug(v ...interface{}) {
	xl.logger.Debug(v...)
}

func (xl *xormLogger) Debugf(f string, v ...interface{}) {
	xl.logger.Debugf(f, v...)
}

func (xl *xormLogger) Info(v ...interface{}) {
	xl.logger.Debug(v...)
}

func (xl *xormLogger) Infof(f string, v ...interface{}) {
	xl.logger.Debugf(f, v...)
}

func (xl *xormLogger) Warn(v ...interface{}) {
	xl.logger.Debug(v...)
}

func (xl *xormLogger) Warnf(f string, v ...interface{}) {
	xl.logger.Errorf(f, v...)
}

func (xl *xormLogger) Error(v ...interface{}) {
	xl.logger.Debug(v...)
}

func (xl *xormLogger) Errorf(f string, v ...interface{}) {
	xl.logger.Errorf(f, v...)
}

func (xl *xormLogger) Level() log.LogLevel {
	return 0
}

func (xl *xormLogger) SetLevel(l log.LogLevel) {
}

func (xl *xormLogger) ShowSQL(b ...bool) {
}

func (xl *xormLogger) IsShowSQL() bool {
	return true
}

// Xorm 初始化Xorm
func Xorm(conf *Config) (*xorm.Engine, error) {
	db, err := xorm.NewEngine(conf.Driver, conf.URI)

	if err != nil {
		return nil, err
	}

	if conf.Debug {
		db.ShowSQL(conf.Debug)
		if conf.Logger != nil {
			logger := log.NewLoggerAdapter(&xormLogger{conf.Logger})
			db.SetLogger(logger)
		}
	}

	return db, nil
}
