# 数据库封装

封装 orm 和数据库驱动

## 初始化

 * 使用 db.Gorm 使用 gorm 初始化
 * 使用 db.Xorm 使用 xorm 初始化

```go
import (
  "github.com/go-eyas/toolkit/db"
  "github.com/go-eyas/toolkit/log"
)

func main() {
  log.Init(&log.Config{})
  var err error
  // gorm
  var db *gorm.DB
  db, err = db.Gorm(db.Config{
    Driver: "mysql",
    URI: "user:password@(127.0.0.1:3306)/mydb",
    Logger: log.SugaredLogger,
  })

  // xorm
  // var db *xorm.Engine
  // db, err = db.Xorm(db.Config{
  //   Driver: "mysql",
  //   URI: "user:password@(127.0.0.1:3306)/mydb",
  //   Logger: log.SugaredLogger,
  // })
  

  if err != nil {
    panic(err)
  }

  defer db.Close()
}
```

## 视图 View
在支持视图的数据库，可使用 `db.GormViewMigrate` 或 `db.XormViewMigrate` 用于创建视图，视图的字段名称映射和 model 一致

视图的模型要实现接口

```go
type ViewModel interface{
	From() string // From 返回 创建视图时的 FROM 部分语句
}
```

```go
type User struct {
  ID int64
  UserName string
  Status byte
}
type Company struct {
  UID int64
  CompanyName string
}

type UserCompany struct {
  *User
  *Company
}

func (UserView) From() string {
  return "FROM users JOIN company ON company.uid = users.uid"
}

db.GormViewMigrate(DB, &UserCompany{})

DB.Model(UserCompany{}).Where("id = ?", 1).Find(&userCompany)
```


## 驱动 

初始化的时候，配置项为 

```go
// Config 数据库配置项
type Config struct {
	Driver string `yaml:"driver" json:"driver" toml:"driver" env:"DB_DRIVER"`
	URI    string `yaml:"uri" json:"uri" toml:"uri" env:"DB_URI"`
	Debug  bool
	Logger Logger
}
```

Driver 的可选项为 

 * mysql
 * postgres
 * mssql: gorm 为 mssql，xorm 为 sqlserver

这些驱动都已提前导入，初始化的时候无需再导入驱动

#### sqlite

因为sqlite驱动是CGO的包，所以默认不导入， 如果要是用sqlite数据库，请按照以下指引

1. 导入驱动
  ```go
  import "github.com/go-eyas/toolkit/db/sqlite"
  ```
2. 安装 Gcc, G++ 编译环境，windows可使用 [TDM-GCC](http://tdm-gcc.tdragon.net/download) ，其他系统的自行解决
3. 使用环境变量启用CGO: `CGO_ENABLED=1`

#### 其他数据库

如果要是用其他数据库，如 oracle，tidb等等，执行查找资料并引入驱动


# godoc

[API 文档](https://gowalker.org/github.com/go-eyas/toolkit/db)