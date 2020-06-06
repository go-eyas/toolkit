package resource

import (
  "github.com/go-eyas/toolkit/db"
  "github.com/jinzhu/gorm"
  "testing"
)

type Article struct {
  ID      int64  `resource:"pk;search:=;order:desc" json:"id"`
  Title   string `resource:"create;update;search:like" json:"title"`
  Content string `resource:"create;update;search:like" json:"text"`
  Status  byte   `resource:"search:=" json:"-"`
}

func testDB() *gorm.DB {
  DB, err := db.Gorm(&db.Config{
    Debug:  true,
    Driver: "mysql",
    URI:    "root:eyas8825345liu@(111.230.219.41:3306)/toolkit_test",
  })
  if err != nil {
    panic(err)
  }
  return DB
}

func TestCreate(t *testing.T) {
  DB := testDB()
  r := NewGormResource(DB, Article{})
  DB.AutoMigrate(&Article{})

  //   create by struct
  err := r.Create(&Article{
    Title:   "测试文章 origin",
    Content: "文章的内容",
  })
  if err != nil {
    panic(err)
  }

  // create by tmp struct
  err = r.Create(&struct{
    Title string
    Content string
  }{
    Title:   "测试文章 tmp struct",
    Content: "文章的内容",
  })
  if err != nil {
    panic(err)
  }

  // create by map
  err = r.Create(map[string]interface{}{
    "title":   "测试文章 map",
    "content": "文章的内容",
  })
  if err != nil {
    panic(err)
  }

  list := []*Article{}
  total, err := r.List(struct{}{}, &list)
  if err != nil {
    panic(err)
  }
  t.Logf("total=%d list=%+v", total, list)


}

// func TestGormResource(t *testing.T) {
//   DB := testDB()
//   r := NewGormResource(DB, Article{})
//   DB.AutoMigrate(&Article{})
//
//   //   create
//   article := &Article{
//     Title:   "测试文章",
//     Content: "文章的内容",
//   }
//   err = r.Create(article)
//   if err != nil {
//     panic(err)
//   }
//
//   err = r.Update(article.ID, map[string]interface{}{
//     "content": "修改后的文章内容",
//   })
//
//   if err != nil {
//     panic(err)
//   }
//
//   err = r.Delete(article.ID)
//
//   if err != nil {
//     panic(err)
//   }
//
// }
