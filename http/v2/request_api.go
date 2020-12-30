package http

import (
	"net/http"
	"time"
)

// 批量设置 HTTP Client 的配置项
type ClientConfig struct {
	TransformRequest  requestMiddlewareHandler // 请求中间件
	TransformResponse responseMiddlewareHandler // 响应中间件
	Headers            map[string]string // 预设 Header
	Cookies            map[string]string // 预设请求 Cookie
	Type              string // 预设请求数据类型, 该配置为空时忽略
	UserAgent         string // 预设 User-Agent, 该配置为空时忽略
	Proxy             string // 预设请求代理，该配置为空时忽略
	BaseURL           string // 预设 url 前缀，叠加
	Query             interface{} // 预设请求查询参数
	Timeout           time.Duration // 请求超时时间，该配置为 0 时忽略
	Retry             int // 重试次数，该配置为 0 时忽略
}

// Config 批量设置 HTTP Client 的配置项
func (c *Client) Config(conf *ClientConfig) *Client {
	cli := c.getSetting()

	if conf.TransformRequest != nil {
		cli = cli.TransformRequest(conf.TransformRequest)
	}

	if conf.TransformResponse != nil {
		cli = cli.TransformResponse(conf.TransformResponse)
	}

	if len(conf.Headers) > 0 {
		for k, v := range conf.Headers {
			cli = cli.Header(k, v)
		}
	}

	if len(conf.Cookies) > 0 {
		for k, v := range conf.Cookies {
			cli = cli.Cookie(k, v)
		}
	}

	if conf.Type != "" {
		cli = cli.Type(conf.Type)
	}
	if conf.UserAgent != "" {
		cli = cli.UserAgent(conf.UserAgent)
	}

	if conf.Proxy != "" {
		cli = cli.Proxy(conf.Proxy)
	}

	if conf.BaseURL != "" {
		cli = cli.BaseURL(conf.BaseURL)
	}

	if conf.Query != nil {
		cli = cli.Query(conf.Query)
	}

	if conf.Timeout > 0 {
		cli = cli.Timeout(conf.Timeout)
	}
	if conf.Retry > 0 {
		cli = cli.Retry(conf.Retry)
	}

	return cli
}

func (c *Client) SetClient(rawClient http.Client) *Client {
	cli := c.getSetting()
	cli.Client = rawClient
	return cli
}

// Header 设置请求 Header
func (c *Client) Header(k string, v interface{}) *Client {
	cli := c.getSetting()
	cli.headers[k] = v
	return cli
}

// TransformRequest 增加请求中间件，可以在请求发起前对整个请求做前置处理，比如修改 body, header, proxy, url, 配置项等等，也可以获取请求的各种数据，如自定义日志，类似于 axios 的 transformRequest
func (c *Client) TransformRequest(h requestMiddlewareHandler) *Client {
	cli := c.getSetting()
	cli.reqMdls = append(cli.reqMdls, h)
	return cli
}

// TransformResponse 增加响应中间件，可以在收到请求后第一时间对请求做处理，如验证 status code，验证 body 数据，甚至重置 body 数据，更改响应等等任何操作，类似于 axios 的 transformResponse
func (c *Client) TransformResponse(h responseMiddlewareHandler) *Client {
	cli := c.getSetting()
	cli.resMdls = append(cli.resMdls, h)
	return cli
}

type ClientMiddleware interface {
	TransformRequest(*Client, *http.Request) *Client
	TransformResponse(*Client, *http.Request, *Response) *Response
}

// Use 应用中间件，实现了 TransformRequest 和 TransformResponse 接口的中间件，如 http.Use(http.AccessLogger())，通常用于成对的请求响应处理
func (c *Client) Use(mdl ClientMiddleware) *Client {
	cli := c.getSetting()
	cli = cli.TransformRequest(mdl.TransformRequest)
	cli = cli.TransformResponse(mdl.TransformResponse)
	return cli
}

// Type 请求提交方式，默认json
func (c *Client) Type(ty string) *Client {
	cli := c.getSetting()
	ct, ok := Types[ty]
	if !ok {
		ct = Types[TypeJSON]
	}
	cli = cli.Header(headerContentTypeKey, ct)
	return cli
}

// UserAgent 设置请求 user-agent
func (c *Client) UserAgent(name string) *Client {
	cli := c.getSetting()
	cli = cli.Header("User-Agent", name)
	return cli
}

// Cookie 设置请求 Cookie
func (c *Client) Cookie(k string, v string) *Client {
	cli := c.getSetting()
	cli.cookies[k] = v
	return cli
}

// Proxy 设置请求代理
func (c *Client) Proxy(url string) *Client {
	cli := c.getSetting()
	cli.proxy = url
	return cli
}

// Query 增加查询参数, 如果设置过多次，将会叠加拼接
func (c *Client) Query(query interface{}) *Client {
	cli := c.getSetting()
	cli.queryArgs = append(cli.queryArgs, query)
	return cli
}

// Timeout 请求超时时间
func (c *Client) Timeout(timeout time.Duration) *Client {
	cli := c.getSetting()
	cli.timeout = timeout
	return cli
}

// BaseURL 设置url前缀
func (c *Client) BaseURL(url string) *Client {
	cli := c.getSetting()
	cli.baseURL += url
	return cli
}

// BaseURL 设置url前缀
func (c *Client) Retry(n int) *Client {
	cli := c.getSetting()
	cli.retryCount = n
	return cli
}

// Head 发起 head 请求
func (c *Client) Head(url string, data ...interface{}) (*Response, error) {
	cli := c.getSetting()
	return cli.DoRequest("HEAD", url, data...)
}

// Get 发起 get 请求
func (c *Client) Get(url string, args ...interface{}) (*Response, error) {
	cli := c.getSetting()
	return cli.DoRequest("GET", url, args...)
}

// Post 发起 post 请求
func (c *Client) Post(url string, args ...interface{}) (*Response, error) {
	cli := c.getSetting()
	return cli.DoRequest("POST", url, args...)
}

// Put 发起 put 请求
func (c *Client) Put(url string, args ...interface{}) (*Response, error) {
	cli := c.getSetting()
	return cli.DoRequest("Put", url, args...)
}

// Del 发起 del 请求
func (c *Client) Del(url string, args ...interface{}) (*Response, error) {
	cli := c.getSetting()
	return cli.DoRequest("DELETE", url, args...)
}

// Patch 发起 patch 请求
func (c *Client) Patch(url string, args ...interface{}) (*Response, error) {
	cli := c.getSetting()
	return cli.DoRequest("Patch", url, args...)
}

// Options 发起 get 请求
func (c *Client) Options(url string, args ...interface{}) (*Response, error) {
	cli := c.getSetting()
	return cli.DoRequest("Options", url, args...)
}
