package http

import "net/http"

// Type 请求提交方式，默认json
func Type(name string) *Request {
	return NewRequest().Type(name)
}

// UserAgent 设置请求 user-agent，默认是 chrome 75.0
func UserAgent(name string) *Request {
	return NewRequest().UserAgent(name)
}

// Cookie 设置请求 Cookie
func Cookie(c *http.Cookie) *Request {
	return NewRequest().Cookie(c)
}

// Header 设置请求 Header
func Header(key, val string) *Request {
	return NewRequest().Header(key, val)
}

// Proxy 设置请求代理
func Proxy(url string) *Request {
	return NewRequest().Proxy(url)
}

// Head 发起 head 请求
func Head(url string, query interface{}) (*Response, error) {
	return NewRequest().Do("HEAD", url, query, nil, nil)
}

// Get 发起 get 请求， query 查询参数
func Get(url string, query interface{}) (*Response, error) {
	return NewRequest().Do("GET", url, query, nil, nil)
}

// Post 发起 post 请求，body 是请求带的参数，可使用json字符串或者结构体
func Post(url string, body interface{}) (*Response, error) {
	return NewRequest().Do("POST", url, nil, body, nil)
}

// Put 发起 put 请求，body 是请求带的参数，可使用json字符串或者结构体
func Put(url string, body interface{}) (*Response, error) {
	return NewRequest().Do("PUT", url, nil, body, nil)
}

// Del 发起 delete 请求，body 是请求带的参数，可使用json字符串或者结构体
func Del(url string, body interface{}) (*Response, error) {
	return NewRequest().Do("DELETE", url, nil, body, nil)
}

// Patch 发起 patch 请求，body 是请求带的参数，可使用json字符串或者结构体
func Patch(url string, body interface{}) (*Response, error) {
	return NewRequest().Do("PATCH", url, nil, body, nil)
}

// Options 发起 options 请求，query 查询参数
func Options(url string, query interface{}) (*Response, error) {
	return NewRequest().Do("OPTIONS", url, query, nil, nil)
}

// PostFile 发起 post 请求上传文件，将使用表单提交，file 是文件地址或者文件流， body 是请求带的参数，可使用json字符串或者结构体
func PostFile(url string, file interface{}, body interface{}) (*Response, error) {
	return NewRequest().Do("PUT", url, nil, body, file)
}

// PutFile 发起 put 请求上传文件，将使用表单提交，file 是文件地址或者文件流， body 是请求带的参数，可使用json字符串或者结构体
func PutFile(url string, file interface{}, body interface{}) (*Response, error) {
	return NewRequest().Do("PUT", url, nil, body, file)
}
