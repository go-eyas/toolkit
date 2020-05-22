package db

import (
  "github.com/go-eyas/toolkit/log"
  "github.com/jinzhu/gorm"
  "testing"
)

type gormModelTest struct {
  gorm.Model
}

func (gormModelTest) TableName() string {
  return "gorm_test"
}

type gormView struct {
  ID int64
}

func (gormView) From() string {
  return `From gorm_test`
}

func TestGorm(t *testing.T) {
  // test init
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
  // test query
  err = db.Raw("SELECT 1 + 1").Row().Scan(&i)
  if err != nil {
    panic(err)
  }
  if i == 2 {
    t.Log("test gorm success")
  }

  // test migrate
  db.AutoMigrate(gormModelTest{})
  list := []*gormModelTest{}
  err = db.Model(gormModelTest{}).Find(&list).Error
  if err != nil {
    panic(err)
  }

  // 	test view
  v := gormView{}
  GormViewMigrate(db, v)
  listView := []*gormView{}
  err = db.Model(v).Find(&listView).Error
  if err != nil {
    panic(err)
  }
}
