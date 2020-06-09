package resource

import (
  "github.com/go-eyas/toolkit/db"
  "github.com/jinzhu/gorm"
)

func New(conf *db.Config, model interface{}) (*Resource, *gorm.DB, error) {
  db, err := db.Gorm(conf)
  if err != nil {
    return nil, nil, err
  }
  r := NewGormResource(db, model)
  return r, db, nil
}