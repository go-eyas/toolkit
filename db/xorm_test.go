package db

import (
	"github.com/go-eyas/toolkit/log"
	"testing"
)

func TestXorm(t *testing.T) {
	db, err := Xorm(&Config{
		Debug:  true,
		Driver: "mysql",
		URI:    "root:123456@(10.0.2.252:3306)/test",
		Logger: log.SugaredLogger,
	})
	if err != nil {
		panic(err)
	}
	i := 0
	_, err = db.SQL("SELECT 1 + 1").Get(&i)
	if err != nil {
		panic(err)
	}

	if i == 2 {
		t.Log("test xorm success")
	}
	type XormTest struct {
		ID int64 `xorm:"id"`
	}
	db.Sync2(XormTest{})
	list := []*XormTest{}
	err = db.Table(XormTest{}).Find(&list)
	if err != nil {
		panic(err)
	}

}
