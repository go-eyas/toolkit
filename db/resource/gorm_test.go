package resource

import (
  "github.com/go-eyas/toolkit/db"
  "testing"
)

type Article struct {
  ID      int64  `resource:"pk;search:none"`
  Title   string `resource:"create;update;search:like"`
  Content string `resource:"create;update;search:like"`
  Status  byte   `resource:"search:="`
}

func TestGormResource(t *testing.T) {
  DB, err := db.Gorm(&db.Config{
    Debug:  true,
    Driver: "mysql",
    URI:    "root:123456@(10.0.2.252:3306)/test",
  })
  if err != nil {
    panic(err)
  }
  r := NewGormResource(DB, Article{})
  DB.AutoMigrate(&Article{})

  //   create
  article := &Article{
    Title:   "测试文章",
    Content: "文章的内容",
  }
  err = r.Create(article)
  if err != nil {
    panic(err)
  }

  err = r.Update(article.ID, map[string]interface{}{
    "content": "修改后的文章内容",
  })

  if err != nil {
    panic(err)
  }

  err = r.Delete(article.ID)

  if err != nil {
    panic(err)
  }

}
