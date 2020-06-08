package db

import (
	"github.com/go-eyas/toolkit/log"
	"os"
	"testing"
)

func TestXorm(t *testing.T) {
	db, err := Xorm(&Config{
		Debug:  true,
		Driver: "mysql",
		URI:    os.Getenv("DB"),
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
