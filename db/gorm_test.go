package db

import (
	"github.com/go-eyas/toolkit/log"
	"github.com/jinzhu/gorm"
	"testing"
)

func TestGorm(t *testing.T) {
	db, err := Gorm(&Config{
		Debug:  true,
		Driver: "mysql",
		URI:    "root:123456@(10.0.2.252:3306)/test",
		Logger: log.SugaredLogger,
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
	type ModelTest struct {
		gorm.Model
	}
	db.AutoMigrate(ModelTest{})
	list := []*ModelTest{}
	err = db.Model(ModelTest{}).Find(&list).Error
	if err != nil {
		panic(err)
	}
}
