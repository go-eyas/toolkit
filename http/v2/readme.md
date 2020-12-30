# HTTP 客户端

封装http 客户端

## 使用

```go
import (
  "fmt"
  "github.com/go-eyas/toolkit/http/v2"
)

func main() {
  h := http.
    TransformRequest(func(client *http.Client, req *nethttp.Request) *http.Client {
        fmt.Printf("HTTP SEND %s %s header=%v\n", req.Method, req.URL, req.Header)
        return client
    }).
    TransformResponse(func (c *http.Client, req *nethttp.Request, resp *http.Response) *http.Response {
        fmt.Printf("HTTP RECV %s %s %d\n", req.Method, req.URL, resp.StatusCode())
        if resp.StatusCode() >= 400 {
            resp.SetBody([]byte("error! error!"))	
        }   
        return resp
    }).
    Type("json").
    Use(http.AccessLogger(logger))
    Header("Authorization", "Bearer xxxxxxxxxxxxxxx").
    Header("x-test", func() string { return "in func string" })
    Header("x-test2", func(cli *http.Client) string { return "in func string2" })
    UserAgent("your custom user-agent").
    Cookie("sid", "sgf2fdas").
    BaseURL("https://api.github.com").
  	Config(&http.Config{
  	    BaseURL: "/api", // 叠加	
    })

  res, err := h.Get("/repos/eyasliu/blog/issues", map[string]string{
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
// get url, 第二个参数可忽略
http.Get("https://api.github.com")

// 带查询参数
http.Get("https://www.google.com", "hl=zh-Hans") // 查询参数可以是字符串
http.Get("https://www.google.com", map[string]string{
	"hl": "zh-Hans",
}) // 可以是map
http.Get("https://www.google.com", struct{
  HL string `json:"hl"`
}{"zh-Hans"}) // 可以是结构体，使用json key作为查询参数的key

// post 请求,第二个参数可忽略
http.Post("https://api.github.com")

// post 带json body参数
http.Post("https://api.github.com", `{"hello": "world"}`) // 可以是字符串
http.Post("https://api.github.com", map[string]interface{}{"hello": "world"}) // 可以是map
http.Post("https://api.github.com", struct{
  Hello string `json:"hello"`
}{"world"}) // 可以是结构体，使用json 序列化字符串

// post 带 查询参数，带json body 参数
http.Query("hl=zh-Hans").Post("https://api.github.com", `{"hello": "world"}`)

// post form表单
http.Type("multipart").Post("https://api.github.com", map[string]interface{}{"hello": "world"})
// post 上传文件，会以表单提交
http.Post("https://api.github.com", "name=@file:./example_file.txt&name=@file:./example_file.txt")

// post 上传多文件，使用file文件流
http.Post("https://api.github.com", "name=@file:./example_file.txt&name=@file:./example_file.txt")

// put, 和post完全一致
http.Put("https://api.github.com")

// delete, 和post完全一致
http.Del("https://api.github.com")

// patch, 和post完全一致
http.Patch("https://api.github.com")

// head, 和get完全一致
http.Head("https://api.github.com")

// options, 和get完全一致
http.Options("https://api.github.com")
```

#### 响应示例

```go
res, err := http.Options("https://api.github.com", nil)

// 错误信息
if err != nil {
  err.Error() // 错误信息
}
res.Error() // 与上面等价

// 响应数据
// 将响应数据转为字符串
var str string = res.String() 

// 将响应数据转为字节
var bt []byte = res.Byte()

// 获取响应状态码
var statusCode = res.Status()

// 获取响应的 header
var headers http.Header = res.Header()

// 获取响应的 cookies
var cookies []*http.Cookie = res.Cookies()

// 与 json 结构体绑定
type ResTest struct {
  Hello string `json:"hello"`
}
rt := &ResTest{}
res.JSON(rt)

// 与 xml 结构体绑定
type ResTest struct {
Hello string `xml:"hello"`
}
rt := &ResTest{}
res.XML(rt)
```

**注意：**

* http的响应状态码 >= 400 时会被视为错误，err 值是 `fmt.Errorf("http status code %d", statusCode)`

#### 链式安全调用

默认是链式安全的，即在链式调用的时候，返回的 `*http.Client` 是个新的实例，不会影响之前链式阶段的配置，如 

```go
cli := http.BaseURL("http://xxx.com")
cli2 := cli.BaseURL("/api") // 需要重新给 cli 赋值才会让其生效

cli.Get("/users") // GET http://xxx.com/users
cli2.Get("/users") // GET http://xxx.com/api/users
```

如果不希望链式安全调用，可以关闭

```go
cli := http.Safe(false) // 关闭后要赋值一次

// 下面的赋值都将生效，链式安全关闭后，赋值和不赋值没有区别
cli.BaseURL("http://xxx.com")
cli = cli.BaseURL("/api") // 可赋值，可也不赋值，没有区别

cli.Get("/users") // GET http://xxx.com/api/users

// 可以后面再进行开启
cli = cli.Safe(true)
cli = cli.BaseURL("/v1") // 开启后如果不赋值将不会生效

```

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
http.TransformRequest(func(client *http.Client, req *nethttp.Request) *http.Client {
  fmt.Printf("HTTP SEND %s %s header=%v\n", req.Method, req.URL, req.Header)
  return client
}).
TransformResponse(func (c *http.Client, req *nethttp.Request, resp *http.Response) *http.Response {
  fmt.Printf("HTTP RECV %s %s %d\n", req.Method, req.URL, resp.StatusCode())
  if resp.StatusCode() >= 400 {
    resp.SetBody([]byte("error! error!"))
  }
  return resp
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

如果是文件上传，会自动设置为 `multipart`，无需手动指定

## godoc

[API 文档](https://gowalker.org/github.com/go-eyas/toolkit/http/v2)