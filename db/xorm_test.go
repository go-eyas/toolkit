package db

import (
	"testing"
)

func TestXorm(t *testing.T) {
	db, err := Xorm(&Config{
		Driver: "mysql",
		URI:    "root:123456@(10.0.3.252:3306)/test",
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

}
