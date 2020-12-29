package http

import "time"

// Header 设置请求 Header
func (r *Request) Header(k string, v interface{}) *Request {
	req := r.getSetting()
	req.headers[k] = v
	return req
}

// TransformRequest 增加请求中间件
func (r *Request) TransformRequest(h requestMiddlewareHandler) *Request {
	req := r.getSetting()
	req.reqMdls = append(req.reqMdls, h)
	return req
}

// TransformResponse 增加响应中间件
func (r *Request) TransformResponse(h responseMiddlewareHandler) *Request {
	req := r.getSetting()
	req.resMdls = append(req.resMdls, h)
	return req
}

// Type 请求提交方式，默认json
func (r *Request) Type(ty string) *Request {
	req := r.getSetting()
	ct, ok := Types[ty]
	if !ok {
		ct = Types[TypeJSON]
	}
	req.Header(headerContentTypeKey, ct)
	return req
}

// UserAgent 设置请求 user-agent，默认是 chrome 75.0
func (r *Request) UserAgent(name string) *Request {
	req := r.getSetting()
	req.Header("User-Agent", name)
	return req
}

// Cookie 设置请求 Cookie
func (r *Request) Cookie(k string, v string) *Request {
	req := r.getSetting()
	req.cookies[k] = v
	return req
}

// Proxy 设置请求代理
func (r *Request) Proxy(url string) *Request {
	req := r.getSetting()
	req.proxy = url
	return req
}

// Query 增加查询参数
func (r *Request) Query(query interface{}) *Request {
	req := r.getSetting()
	req.queryArgs = append(req.queryArgs, query)
	return req
}

// Timeout 请求超时时间
func (r *Request) Timeout(timeout time.Duration) *Request {
	req := r.getSetting()
	req.timeout = timeout
	return req
}

// BaseURL 设置url前缀
func (r *Request) BaseURL(url string) *Request {
	req := r.getSetting()
	req.baseURL += url
	return req
}

// Head 发起 head 请求
func (r *Request) Head(url string, data ...interface{}) (*Response, error) {
	req := r.getSetting()
	return req.DoRequest("HEAD", url, data...)
}

// Get 发起 get 请求
func (r *Request) Get(url string, args ...interface{}) (*Response, error) {
	req := r.getSetting()
	return req.DoRequest("GET", url, args...)
}

// Post 发起 post 请求
func (r *Request) Post(url string, args ...interface{}) (*Response, error) {
	req := r.getSetting()
	return req.DoRequest("POST", url, args...)
}

// Put 发起 put 请求
func (r *Request) Put(url string, args ...interface{}) (*Response, error) {
	req := r.getSetting()
	return req.DoRequest("Put", url, args...)
}

// Del 发起 del 请求
func (r *Request) Del(url string, args ...interface{}) (*Response, error) {
	req := r.getSetting()
	return req.DoRequest("DELETE", url, args...)
}

// Patch 发起 patch 请求
func (r *Request) Patch(url string, args ...interface{}) (*Response, error) {
	req := r.getSetting()
	return req.DoRequest("Patch", url, args...)
}

// Options 发起 get 请求
func (r *Request) Options(url string, args ...interface{}) (*Response, error) {
	req := r.getSetting()
	return req.DoRequest("Options", url, args...)
}