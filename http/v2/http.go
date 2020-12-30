package http

import "time"

var defClient = New().Type("json")


// Type 请求提交方式，默认json
func Type(name string) *Client {
	return defClient.Type(name)
}

// UserAgent 设置请求 user-agent，默认是 chrome 75.0
func UserAgent(name string) *Client {
	return defClient.UserAgent(name)
}

// Cookie 设置请求 Cookie
func Cookie(k string, v string) *Client {
	return defClient.Cookie(k, v)
}

// Header 设置请求 Header
func Header(key, val string) *Client {
	return defClient.Header(key, val)
}

// Proxy 设置请求代理
func Proxy(url string) *Client {
	return defClient.Proxy(url)
}

// Query 设置请求代理
func Query(query interface{}) *Client {
	return defClient.Query(query)
}

// Timeout 设置请求代理
func Timeout(timeout time.Duration) *Client {
	return defClient.Timeout(timeout)
}

// UseRequest 增加请求中间件
func Use(mdl ClientMiddleware) *Client {
	return defClient.Use(mdl)
}

// UseResponse 增加响应中间件
func UseResponse(mdl responseMiddlewareHandler) *Client {
	return defClient.TransformResponse(mdl)
}

// BaseURL 设置url前缀
func BaseURL(url string) *Client {
	return defClient.BaseURL(url)
}

// BaseURL 设置url前缀
func Retry(n int) *Client {
	return defClient.Retry(n)
}

// Head 发起 head 请求
func Head(url string, data ...interface{}) (*Response, error) {
	return defClient.Head(url, data...)
}

// Get 发起 get 请求， query 查询参数
func Get(url string, data ...interface{}) (*Response, error) {
	return defClient.Get(url, data...)
}

// Post 发起 post 请求，body 是请求带的参数，可使用json字符串或者结构体
func Post(url string, data ...interface{}) (*Response, error) {
	return defClient.Post(url, data...)
}

// Put 发起 put 请求，body 是请求带的参数，可使用json字符串或者结构体
func Put(url string, data ...interface{}) (*Response, error) {
	return defClient.Put(url, data...)
}

// Del 发起 delete 请求，body 是请求带的参数，可使用json字符串或者结构体
func Del(url string, data ...interface{}) (*Response, error) {
	return defClient.Del(url, data...)
}

// Patch 发起 patch 请求，body 是请求带的参数，可使用json字符串或者结构体
func Patch(url string, data ...interface{}) (*Response, error) {
	return defClient.Patch(url, data...)
}

// Options 发起 options 请求，query 查询参数
func Options(url string, data ...interface{}) (*Response, error) {
	return defClient.Options(url, data...)
}
