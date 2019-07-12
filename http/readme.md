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

// 与结构体绑定
type ResTest struct {
  Hello string `json:"hello"`
}
rt := &ResTest{}
res.JSON(rt)
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

```
package http // import "github.com/go-eyas/toolkit/http123"      


TYPES

type Request struct {
        SuperAgent *gorequest.SuperAgent
        Req        *gorequest.Request

        // Has unexported fields.
}
    Request 请求结构

func Cookie(c *http.Cookie) *Request
    Cookie 设置请求 Cookie

func Header(key, val string) *Request
    Header 设置请求 Header

func New() *Request
    New 新建请求对象，默认useragent 为 chrome 75.0, 数据类型 json

func Proxy(url string) *Request
    Proxy 设置请求代理

func Query(query interface{}) *Request
    Query 设置请求代理

func Timeout(timeout time.Duration) *Request
    Timeout 设置请求代理

func Type(name string) *Request
    Type 请求提交方式，默认json

func UseRequest(mdl requestMiddlewareHandler) *Request
    UseRequest 增加请求中间件

func UseResponse(mdl responseMidlewareHandler) *Request
    UseResponse 增加响应中间件

func UserAgent(name string) *Request
    UserAgent 设置请求 user-agent，默认是 chrome 75.0

func (r *Request) Cookie(c *http.Cookie) *Request
    Cookie 设置请求 Cookie

func (r *Request) Del(url string, body interface{}) (*Response, error)
    Del 发起 delete 请求，body 是请求带的参数，可使用json字符串或者结构体

func (r *Request) Do(method, url string, args ...interface{}) (*Response, error)
    Do 发出请求，method 请求方法，url 请求地址， query 查询参数，body 请求数据，file 文件对象/地址

func (r *Request) Get(url string, query interface{}) (*Response, error)
    Get 发起 get 请求， query 查询参数

func (r *Request) Head(url string, query interface{}) (*Response, error)
    Head 发起 head 请求

func (r *Request) Header(key, val string) *Request
    Header 设置请求 Header

func (r *Request) Options(url string, query interface{}) (*Response, error)
    Options 发起 options 请求，query 查询参数

func (r *Request) Patch(url string, body interface{}) (*Response, error)
    Patch 发起 patch 请求，body 是请求带的参数，可使用json字符串或者结构体

func (r *Request) Post(url string, body interface{}) (*Response, error)
    Post 发起 post 请求，body 是请求带的参数，可使用json字符串或者结构体

func (r *Request) PostFile(url string, file interface{}, body interface{}) (*Response, error)
    PostFile 发起 post 请求上传文件，将使用表单提交，file 是文件地址或者文件流， body
    是请求带的参数，可使用json字符串或者结构体

func (r *Request) Proxy(url string) *Request
    Proxy 设置请求代理

func (r *Request) Put(url string, body interface{}) (*Response, error)
    Put 发起 put 请求，body 是请求带的参数，可使用json字符串或者结构体

func (r *Request) PutFile(url string, file interface{}, body interface{}) (*Response, error)
    PutFile 发起 put 请求上传文件，将使用表单提交，file 是文件地址或者文件流， body 是请求带的参数，可使用json字符串或者结构体

func (r *Request) Query(query interface{}) *Request
    Query 增加查询参数

func (r *Request) Timeout(timeout time.Duration) *Request
    Timeout 请求超时时间

func (r *Request) Type(name string) *Request
    Type 请求提交方式，默认json

func (r *Request) UseRequest(mdl requestMiddlewareHandler) *Request
    UseRequest 增加请求中间件

func (r *Request) UseResponse(mdl responseMidlewareHandler) *Request
    UseResponse 增加响应中间件

func (r *Request) UserAgent(name string) *Request
    UserAgent 设置请求 user-agent，默认是 chrome 75.0

type Response struct {
        Request *Request
        Raw     *http.Response
        Body    []byte
        Errs    ResponseError
}
    Response 回应对象

func Del(url string, body interface{}) (*Response, error)
    Del 发起 delete 请求，body 是请求带的参数，可使用json字符串或者结构体

func Get(url string, query interface{}) (*Response, error)
    Get 发起 get 请求， query 查询参数

func Head(url string, query interface{}) (*Response, error)
    Head 发起 head 请求

func NewResponse() *Response
    NewResponse 新建回应对象

func Options(url string, query interface{}) (*Response, error)
    Options 发起 options 请求，query 查询参数

func Patch(url string, body interface{}) (*Response, error)
    Patch 发起 patch 请求，body 是请求带的参数，可使用json字符串或者结构体

func Post(url string, body interface{}) (*Response, error)
    Post 发起 post 请求，body 是请求带的参数，可使用json字符串或者结构体

func PostFile(url string, file interface{}, body interface{}) (*Response, error)
    PostFile 发起 post 请求上传文件，将使用表单提交，file 是文件地址或者文件流， body
    是请求带的参数，可使用json字符串或者结构体

func Put(url string, body interface{}) (*Response, error)
    Put 发起 put 请求，body 是请求带的参数，可使用json字符串或者结构体

func PutFile(url string, file interface{}, body interface{}) (*Response, error)
    PutFile 发起 put 请求上传文件，将使用表单提交，file 是文件地址或者文件流， body 是请求带的参数，可使用json字符串或者结构体

func (r *Response) Byte() []byte
    Byte 获取响应字节

func (r *Response) Cookies() []*http.Cookie
    Cookies 获取响应 cookie

func (r *Response) Err() error
    Err 获取响应错误

func (r *Response) Header() http.Header
    Header 获取响应header

func (r *Response) IsError() bool
    IsError 是否响应错误

func (r *Response) JSON(v interface{}) error
    JSON 根据json绑定结构体

func (r *Response) Status() int
    Status 获取响应状态码

func (r *Response) String() string
    String 获取响应字符串

type ResponseError []error
    ResponseError 响应错误对象

func (e ResponseError) Add(err error) ResponseError
    Add 增加错误

func (e ResponseError) Error() string
    Error 实现 error 接口

func (e ResponseError) HasErr() bool
    HasErr 是否有错误
```