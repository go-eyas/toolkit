package db

import (
	"testing"
)

func TestGorm(t *testing.T) {
	db, err := Gorm(&Config{
		Driver: "mysql",
		URI:    "root:123456@(10.0.3.252:3306)/test",
	})
	if err != nil {
		panic(err)
	}
	i := 0
	err = db.Raw("SELECT 1 + 1").Row().Scan(&i)
	if err != nil {
		panic(err)
	}
	if i == 2 {
		t.Log("test gorm success")
	}
}
