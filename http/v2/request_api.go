package http

import (
	"net/http"
	"time"
)

// Header 设置请求 Header
func (c *Client) Header(k string, v interface{}) *Client {
	cli := c.getSetting()
	cli.headers[k] = v
	return cli
}

// TransformRequest 增加请求中间件
func (c *Client) TransformRequest(h requestMiddlewareHandler) *Client {
	cli := c.getSetting()
	cli.reqMdls = append(cli.reqMdls, h)
	return cli
}

// TransformResponse 增加响应中间件
func (c *Client) TransformResponse(h responseMiddlewareHandler) *Client {
	cli := c.getSetting()
	cli.resMdls = append(cli.resMdls, h)
	return cli
}

type ClientMiddleware interface {
	TransformRequest(*Client, *http.Request) *Client
	TransformResponse(*Client, *http.Request, *Response) *Response
}

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
	cli.Header(headerContentTypeKey, ct)
	return cli
}

// UserAgent 设置请求 user-agent，默认是 chrome 75.0
func (c *Client) UserAgent(name string) *Client {
	cli := c.getSetting()
	cli.Header("User-Agent", name)
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

// Query 增加查询参数
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