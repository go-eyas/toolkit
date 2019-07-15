# HTTP 客户端

封装http 客户端

## 使用

```go
import (
  "fmt"
  "github.com/go-eyas/toolkit/http"
)

func main() {
  h := http.Header("Authorization", "Bearer xxxxxxxxxxxxxxx").UserAgent("your custom user-agent").Cookie()

  res, err := h.Get("https://api.github.com/repos/eyasliu/blog/issues", map[string]string{
    "per_page": 1,
  })

  // 获取字符串
  fmt.Printf("print string: %s\n", res.String())
  // 获取字节
  fmt.Printf("print bytes: %v", res.Byte())

  // 绑定结构体
  s := []struct {
		URL   string `json:"url"`
		Title string `json:"title"`
  }{}
  res.JSON(&s)
  fmt.Printf("print Struct: %v", s)

  // 使用代理
  res, err := http.Proxy("http://127.0.0.1:1080").Get("https://www.google.com", map[string]string{
		"hl": "zh-Hans",
  })
  fmt.Printf("google html: %s", res.String())
}
```

## 使用指南

#### 请求示例

```go
// get url
http.Get("https://api.github.com", nil)

// 带查询参数
http.Get("https://www.google.com", "hl=zh-Hans") // 查询参数可以是字符串
http.Get("https://www.google.com", map[string]string{
	"hl": "zh-Hans",
}) // 可以是map
http.Get("https://www.google.com", struct{
  HL string `json:"hl"`
}{"zh-Hans"}) // 可以是结构体，使用json key作为查询参数的key

// post 请求
http.Post("https://api.github.com", nil)

// post 带json参数
http.Post("https://api.github.com", `{"hello": "world"}`) // 可以是字符串
http.Post("https://api.github.com", map[string]interface{}{"hello": "world"}) // 可以是map
http.Post("https://api.github.com", struct{
  Hello string `json:"hello"`
}{"world"}) // 可以是结构体，使用json 序列化字符串

// post 带 查询参数，带json参数
http.Query("hl=zh-Hans").Post("https://api.github.com", `{"hello": "world"}`)

// post form表单
http.Type("multipart").Post("https://api.github.com", map[string]interface{}{"hello": "world"})
// post 上传文件，会以表单提交
http.PostFile("https://api.github.com", "./example_file.txt", map[string]interface{}{"hello": "world"})

// post 上传文件，使用file文件流
file, _ := ioutil.ReadFile("./example_file.txt")
file, _ := os.Open("./example_file.txt")
http.PostFile("https://api.github.com", file, map[string]interface{}{"hello": "world"})

// put, 和post完全一致
http.Put("https://api.github.com", nil)

// delete, 和post完全一致
http.Del("https://api.github.com", nil)

// patch, 和post完全一致
http.Patch("https://api.github.com", nil)

// head, 和get完全一致
http.Head("https://api.github.com", nil)

// options, 和get完全一致
http.Options("https://api.github.com", nil)
```

#### 响应示例

```go
res, err := http.Options("https://api.github.com", nil)

// 错误信息
if err != nil {
  err.Error() // 错误信息
}
res.Err().Error() // 与上面等价

// 响应数据
// 将响应数据转为字符串
var str string = res.String() 

// 将响应数据转为字节
var bt []byte = res.Byte()

// 获取响应状态码
var statusCode = res.Status()

// 获取响应的 header
var http.Header = res.Header()

// 获取响应的 cookies
var []*http.Cookie = res.Cookies()

// 与结构体绑定
type ResTest struct {
  Hello string `json:"hello"`
}
rt := &ResTest{}
res.JSON(rt)
```

**注意：**

 * http的响应状态码 >= 400 时会被视为错误，err 值是 `fmt.Errorf("http response status code %d", statusCode)`

#### 提前设置通用项

```go
h := http.Header("Authorization", "Bearer xxxxxxxxxxxxxxx"). // 设置header
    UserAgent("your custom user-agent"). // 设置 useragent
    Timeout(10 * time.Second). // 设置请求超时时间
    Query("lang=zh_ch"). // 设置查询参数
    Proxy("http://127.0.0.1:1080") // 设置代理

h.Get("xxxx", nil)
```

#### 中间件支持

可以增加请求中间件和响应中间件，用于在请求或响应中改变内部操作

```go
http.UseRequest(func(req *http.Request) *http.Request {
    fmt.Printf("http 发送 %s %s\n", req.SuperAgent.Method, req.SuperAgent.Url)
    return req
}).UseResponse(func(req *http.Request, res *http.Response) *http.Response {
    fmt.Printf("http 接收 %s %s\n", req.SuperAgent.Method, req.SuperAgent.Url)
    return res
})
```

#### 代理设置

默认会获取环境变量 `http_proxy` 的值使用代理，但是可以手动指定

```go
http.Proxy("http://127.0.0.1:1080").Get("https://www.google.com", map[string]string{
	"hl": "zh-Hans",
})

// 临时取消代理
http.Proxy("").Get("https://www.google.com", map[string]string{
	"hl": "zh-Hans",
})
```

#### 提交方式

也就是 `Type(t string)` 函数支持的值

```
"text/html" uses "html"
"application/json" uses "json"
"application/xml" uses "xml"
"text/plain" uses "text"
"application/x-www-form-urlencoded" uses "urlencoded", "form" or "form-data"
```

如果是文件上传，则应该设置为 `multipart`

## godoc

[API 文档](https://gowalker.org/github.com/go-eyas/toolkit/http)