# Resource 资源自动维护、检索

**not finish yet 还没完成**

基于 Gorm API，根据 Restful API 设计对资源进行维护并检索

简单说就是：自动 curd

# Usage

先来看个栗子

```go
package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/go-eyas/toolkit/db"
  "github.com/go-eyas/toolkit/db/resource"
)


type Article struct {
  ID      int64  `resource:"pk;search:none"`
  Title   string `resource:"create;update;search:like"`
  Content string `resource:"create;update;search:like"`
  Status  byte   `resource:"search:="`
}

func main() {
  DB, err := db.Gorm(&db.Config{URI: "root:123456@(127.0.0.1:3306)/test"})
  r := resource.NewGormResource(DB, &Article{})
  
  /******* create resource ********/
  err = r.Create(&Article{Title: "the title", Content: "the content"}) // 使用原本类型结构体

  // 使用临时结构体
  err = r.Create(&struct{
    Title string
    Content string
  }{Title: "the title", Content: "the content"})

  // 使用 map
  err = r.Create(map[string]interface{}{"title": "the title", "content": "the content"})


  /******* update resource ********/
  err = r.Update(1, &Article{Title: "the title", Content: "the content"})
  err = r.Update(1, &struct{
    Title string
    Content string
  }{Title: "the title", Content: "the content"})
  err = r.Update(1, map[string]interface{}{"title": "the title", "content": "the content"})


  /******* delete resource ********/
  err = r.Delete(1)

  /******* get one resource ********/
  err = r.Detail(1)

  /******* list resource ********/
  list := []*Article{}

  total, err := r.List(&Article{Title: "the title", Content: "the content"}, &list)

  total, err := r.List(&struct{
    Title string
    Content string
    Status byte
  }{Status: 1}, &list)

  total, err := r.List(map[string]interface{}{
    "status": 1,
  }, &list)


  /******* 绑定路由 ********/
  // 使用 net/http 
  http.Handle("/articles", r.HTTPHandler)
  http.Handle("/articles/:id", r.HTTPHandler)

  // 使用 gin
  router := gin.Default()
  router.GET("/articles", gin.WrapF(r.HTTPListHandler))  
  router.GET("/articles/:id", gin.WrapF(r.HTTPDetailHandler))  
  router.POST("/articles", gin.WrapF(r.HTTPCreateHandler))  
  router.PUT("/articles/:id", gin.WrapF(r.HTTPUpdateHandler))  
  router.DELETE("/articles/:id", gin.WrapF(r.HTTPDeleteHandler))  

    


}
```