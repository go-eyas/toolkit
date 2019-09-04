package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

type requestMiddlewareHandler func(*Request) *Request
type responseMidlewareHandler func(*Request, *Response) *Response

func newRaw() Request {
	return Request{
		SuperAgent: gorequest.New(),
		reqMdls:    []requestMiddlewareHandler{},
		resMdls:    []responseMidlewareHandler{},
	}
}

// New 新建请求对象，默认数据类型 json
func New() Request {
	r := newRaw()
	r = r.Type("json")

	return r
}

// Request 请求结构
type Request struct {
	SuperAgent  *gorequest.SuperAgent
	querys      []interface{}
	reqMdls     []requestMiddlewareHandler
	resMdls     []responseMidlewareHandler
	cookies     []*http.Cookie
	baseURL     string
	contentType string
	proxy       string
	timeout     time.Duration
}

func (r Request) Clone() Request {
	req := newRaw()
	req.baseURL = r.baseURL
	req.proxy = r.proxy
	req.contentType = r.contentType

	// query
	for _, query := range r.querys {
		req.querys = append(req.querys, query)
	}

	for _, cookie := range r.SuperAgent.Cookies {
		req.SuperAgent.Cookies = append(req.SuperAgent.Cookies, cookie)
	}

	// req mdl
	for _, mdl := range r.reqMdls {
		req.reqMdls = append(req.reqMdls, mdl)
	}

	// res mdl
	for _, mdl := range r.resMdls {
		req.resMdls = append(req.resMdls, mdl)
	}

	// headers
	for k, v := range r.SuperAgent.Header {
		req.SuperAgent.Header[k] = v
	}

	return req
}

// Type 请求提交方式，默认json
func (r Request) Type(name string) Request {
	req := r.Clone()
	req.contentType = name
	return req
}

// UserAgent 设置请求 user-agent，默认是 chrome 75.0
func (r Request) UserAgent(name string) Request {
	req := r.Clone()
	req.SuperAgent = req.SuperAgent.Set("User-Agent", name)
	return req
}

// Cookie 设置请求 Cookie
func (r Request) Cookie(c *http.Cookie) Request {
	req := r.Clone()
	req.cookies = append(req.cookies, c)
	return req
}

// Header 设置请求 Header
func (r Request) Header(key, val string) Request {
	req := r.Clone()
	req.SuperAgent = req.SuperAgent.Set(key, val)
	return req
}

// Proxy 设置请求代理
func (r Request) Proxy(url string) Request {
	req := r.Clone()
	req.proxy = url
	return req
}

// Query 增加查询参数
func (r Request) Query(query interface{}) Request {
	req := r.Clone()
	req.querys = append(req.querys, query)
	return req
}

// Timeout 请求超时时间
func (r Request) Timeout(timeout time.Duration) Request {
	req := r.Clone()
	req.timeout = timeout
	return req
}

// UseRequest 增加请求中间件
func (r Request) UseRequest(mdl requestMiddlewareHandler) Request {
	req := r.Clone()
	req.reqMdls = append(req.reqMdls, mdl)
	return req
}

// UseResponse 增加响应中间件
func (r Request) UseResponse(mdl responseMidlewareHandler) Request {
	req := r.Clone()
	req.resMdls = append(req.resMdls, mdl)
	return req
}

// BaseURL 设置url前缀
func (r Request) BaseURL(url string) Request {
	req := r.Clone()
	req.baseURL += url
	return req
}

// Do 发出请求，method 请求方法，url 请求地址， query 查询参数，body 请求数据，file 文件对象/地址
func (r Request) Do(method, url string, query, body, file interface{}) (*Response, error) {
	r = r.Clone()
	// set mthod url
	if method == "" || url == "" {
		return &Response{
			Request: &r,
			Raw:     nil,
			Body:    []byte{},
			Errs:    []error{errors.New("url is empty")},
		}, fmt.Errorf("http url can't empty")
	}
	// r.SuperAgent = r.SuperAgent.CustomMethod(method, r.baseURL+url)
	r.SuperAgent.Method = strings.ToUpper(method)
	r.SuperAgent.Url = r.baseURL + url
	r.SuperAgent.Errors = nil

	if r.contentType != "" {
		r.SuperAgent = r.SuperAgent.Type(r.contentType)
	}
	if r.timeout > 0 {
		r.SuperAgent = r.SuperAgent.Timeout(r.timeout)
	}

	if r.proxy != "" {
		r.SuperAgent = r.SuperAgent.Proxy(r.proxy)
	}

	// set query string
	if query != nil {
		r.SuperAgent = r.SuperAgent.Query(query)
	}
	for _, q := range r.querys {
		r.SuperAgent = r.SuperAgent.Query(q)
	}

	// set body
	if body != nil {
		r.SuperAgent = r.SuperAgent.Send(body)
	}

	if file != nil {
		r.Type("multipart")
		r.SuperAgent = r.SuperAgent.SendFile(file)
	}

	// 执行请求中间件
	for _, mdl := range r.reqMdls {
		r1 := mdl(&r)
		r = *r1
	}

	res, resBody, errs := r.SuperAgent.EndBytes()

	response := &Response{
		Request: &r,
		Raw:     res,
		Body:    resBody,
		Errs:    errs,
	}

	// 执行响应中间件
	for _, mdl := range r.resMdls {
		response = mdl(&r, response)
	}

	statusCode := response.Status()
	if statusCode >= 400 {
		response.Errs = response.Errs.Add(fmt.Errorf("http response status code %d", statusCode))
	}

	return response, response.Err()
}

// Head 发起 head 请求
func (r Request) Head(url string, query interface{}) (*Response, error) {
	return r.Do("HEAD", url, query, nil, nil)
}

// Get 发起 get 请求， query 查询参数
func (r Request) Get(url string, query interface{}) (*Response, error) {
	return r.Do("GET", url, query, nil, nil)
}

// Post 发起 post 请求，body 是请求带的参数，可使用json字符串或者结构体
func (r Request) Post(url string, body interface{}) (*Response, error) {
	return r.Do("POST", url, nil, body, nil)
}

// Put 发起 put 请求，body 是请求带的参数，可使用json字符串或者结构体
func (r Request) Put(url string, body interface{}) (*Response, error) {
	return r.Do("PUT", url, nil, body, nil)
}

// Del 发起 delete 请求，body 是请求带的参数，可使用json字符串或者结构体
func (r Request) Del(url string, body interface{}) (*Response, error) {
	return r.Do("DELETE", url, nil, body, nil)
}

// Patch 发起 patch 请求，body 是请求带的参数，可使用json字符串或者结构体
func (r Request) Patch(url string, body interface{}) (*Response, error) {
	return r.Do("PATCH", url, nil, body, nil)
}

// Options 发起 options 请求，query 查询参数
func (r Request) Options(url string, query interface{}) (*Response, error) {
	return r.Do("OPTIONS", url, query, nil, nil)
}

// PostFile 发起 post 请求上传文件，将使用表单提交，file 是文件地址或者文件流， body 是请求带的参数，可使用json字符串或者结构体
func (r Request) PostFile(url string, file interface{}, body interface{}) (*Response, error) {
	return r.Do("PUT", url, nil, body, file)
}

// PutFile 发起 put 请求上传文件，将使用表单提交，file 是文件地址或者文件流， body 是请求带的参数，可使用json字符串或者结构体
func (r Request) PutFile(url string, file interface{}, body interface{}) (*Response, error) {
	return r.Do("PUT", url, nil, body, file)
}
