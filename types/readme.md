# 黑魔法类型

一些带有特殊用途的类型

### JSONString

json 字符串，该类型在转成json字符串的时候，会自动转成object类型，所以要求该类型的值是一个合法的json字符串

使用场景：扩展字段

```go
import "github.com/go-eyas/toolkit/types"

var str = types.JSONString(`{"demo": true, "num": 123}`)

data := struct {
  S JSONString
}{str}

json.Marshal(data) // {"S":{"demo":true,"num":123}}

data2 := struct{
  Demo bool
  Num int
}{}
str.JSON(&data2) // 也可以直接 json 反序列化
```

### Time

时间类型 `time.Time` 的别名，该类型在转成 json 字符串的时候，会把时间格式化成这种格式 2006-01-02 15:04:05

结合 gorm 使用，存在数据库的是时间类型，转到接口的是上述时间格式