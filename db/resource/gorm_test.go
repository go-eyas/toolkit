package resource

import (
  "github.com/go-eyas/toolkit/db"
  "github.com/jinzhu/gorm"
  "os"
  "testing"
)

type Article struct {
  ID      int64  `resource:"pk;search:=;order:desc" json:"id"`
  Title   string `resource:"create;update;search:like" json:"title"`
  Content string `resource:"create;update;search:like" json:"text"`
  Status  byte   `resource:"search:=" json:"-"`
}

var dbConfig = &db.Config{
  Debug:  true,
  Driver: "mysql",
  URI:    os.Getenv("DB"),
}

func testDB() *gorm.DB {
  DB, err := db.Gorm(dbConfig)
  if err != nil {
    panic(err)
  }
  return DB
}

func TestCreate(t *testing.T) {
  r, DB, err := New(dbConfig, Article{})
  if err != nil {
    panic(err)
  }
  DB.AutoMigrate(&Article{})

  //   create by struct
  err = r.Create(&Article{
    Title:   "测试文章 origin",
    Content: "文章的内容",
  })
  if err != nil {
    panic(err)
  }

  // create by tmp struct
  err = r.Create(&struct {
    Title   string
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
  total, err := r.List(&list)
  if err != nil {
    panic(err)
  }
  t.Logf("total=%d list=%+v", total, list)
}

func TestUpdate(t *testing.T) {
  DB := testDB()
  r := NewGormResource(DB, Article{})
  DB.AutoMigrate(&Article{})

  testModel := &Article{
    Title:   "测试文章 origin",
    Content: "文章的内容",
    Status:  1,
  }
  err := r.Create(testModel)
  if err != nil {
    panic(err)
  }
  err = r.Update(testModel.ID, &Article{Status: 2})
  if err != nil {
    panic(err)
  }

  err = r.Update(testModel.ID, &struct{ Status int }{Status: 3})
  if err != nil {
    panic(err)
  }

  err = r.Update(testModel.ID, map[string]byte{"status": 0})
  if err != nil {
    panic(err)
  }
}

func TestDelete(t *testing.T) {
  DB := testDB()
  r := NewGormResource(DB, Article{})
  DB.AutoMigrate(&Article{})

  testModel := &Article{
    Title:   "测试文章 origin",
    Content: "文章的内容",
    Status:  1,
  }
  err := r.Create(testModel)
  if err != nil {
    panic(err)
  }

  err = r.Delete(testModel.ID)
  if err != nil {
    panic(err)
  }
}

func TestDetail(t *testing.T) {
  DB := testDB()
  r := NewGormResource(DB, Article{})
  DB.AutoMigrate(&Article{})

  testModel := &Article{
    Title:   "测试文章 origin",
    Content: "文章的内容",
    Status:  1,
  }
  err := r.Create(testModel)
  if err != nil {
    panic(err)
  }

  dest := &Article{}
  err = r.Detail(testModel.ID, dest)
  if err != nil {
    panic(err)
  }
  t.Logf("dest: %+v", dest)
}

func TestList(t *testing.T) {
  DB := testDB()
  r := NewGormResource(DB, Article{})
  DB.AutoMigrate(&Article{})

  list := []*Article{}
  total, err := r.List(&list)
  if err != nil {
    panic(err)
  } else {
    t.Logf("total=%d  list=%+v", total, list)
  }

  list = []*Article{}
  total, err = r.List(&list, &Article{Status: 0})
  if err != nil {
    panic(err)
  } else {
    t.Logf("total=%d  list=%+v", total, list)
  }

  list = []*Article{}
  total, err = r.List(&list, &Article{Title: "测试"})
  if err != nil {
    panic(err)
  } else {
    t.Logf("total=%d  list=%+v", total, list)
  }

  list = []*Article{}
  total, err = r.List(&list, &struct{ Title string }{Title: "测试"})
  if err != nil {
    panic(err)
  } else {
    t.Logf("total=%d  list=%+v", total, list)
  }

  list = []*Article{}
  total, err = r.List(&list, map[string]interface{}{"title": "测试", "status": 0})
  if err != nil {
    panic(err)
  } else {
    t.Logf("total=%d  list=%+v", total, list)
  }

  list = []*Article{}
  total, err = r.List(&list, nil, []string{"id DESC", "status ASC"})
  if err != nil {
    panic(err)
  } else {
    t.Logf("total=%d  list=%+v", total, list)
  }

  list = []*Article{}
  total, err = r.List(&list, nil, map[string]interface{}{
    "id":       "desc",
    "status":   "Asc",
    "intvalue": 11,
  })
  if err != nil {
    panic(err)
  } else {
    t.Logf("total=%d  list=%+v", total, list)
  }

  list = []*Article{}
  total, err = r.List(&list, nil, map[string]string{
    "id":     "desc",
    "status": "Asc",
  })
  if err != nil {
    panic(err)
  } else {
    t.Logf("total=%d  list=%+v", total, list)
  }

  list = []*Article{}
  total, err = r.List(&list, map[string]int{
    "offset": 10,
    "limit":  10,
  }, map[string]string{
    "id":     "desc",
    "status": "Asc",
  })
  if err != nil {
    panic(err)
  } else {
    t.Logf("total=%d  list=%+v", total, list)
  }

}

// test view crud

type UserBase struct {
  ID   int64
  Name string
}
type Profile struct {
  UserID int64
  Avatar string
  Desc   string
}

type UserProfile struct {
  UID    int64 `resource:"pk" gorm:"column:id"`
  Name   string
  Avatar string
  Desc   string
}

func (UserProfile) From() string {
  return `From user_bases LEFT JOIN profiles ON profiles.user_id=user_bases.id`
}

func TestViewCrud(t *testing.T) {
  DB := testDB()
  DB.AutoMigrate(&UserBase{}, &Profile{})
  db.GormViewMigrate(DB, &UserProfile{})
  r := NewGormResource(DB, UserProfile{})

  base := &UserBase{Name: "test User"}
  DB.Model(base).Create(base)
  DB.Model(&Profile{}).Create(&Profile{UserID: base.ID, Avatar: "/test/avatar.png", Desc: "Hello test"})

  userProfile := &UserProfile{}
  err := r.Detail(base.ID, userProfile)
  if err != nil {
    panic(err)
  } else {
    t.Logf("data=%+v", userProfile)
  }

  // err = r.Update(userProfile.UID, map[string]string{
  //   "name": "resource user",
  //   "Desc": "Hello resource",
  // })
  // if err != nil {
  //   panic(err)
  // }

  list := []*UserProfile{}
  total, err := r.List(&list)
  if err != nil {
    panic(err)
  } else {
    t.Logf("total=%d  list=%+v", total, list)
  }

  // err = r.Delete(userProfile.UID)
  // if err != nil {
  //   panic(err)
  // }
}
