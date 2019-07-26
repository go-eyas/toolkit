# 给 gin 定制的工具函数

## response

使接口的回应固定格式为 

```json
{
  "status": 0,
  "msg": "ok",
  "data": {}
}
```

使用 

```go
import (
  "github.com/gin-gonic/gin"
  "github.com/go-eyas/toolkit/gin/util"
)

func HelloHandler(c *gin.Context) {
  util.R(c).OK(gin.H{
    "hello": "world",
  })
  // 将会响应为
  // {
  //   "status": 0,
  //   "msg": "ok",
  //   "data": {
  //     "hello": "world",
  //   }
  // }
}
```