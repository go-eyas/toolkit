package http

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewResponse 新建回应对象
func NewResponse() *Response {
	return &Response{}
}

// ResponseError 响应错误对象
type ResponseError []error

// Error 实现 error 接口
func (e ResponseError) Error() string {
	errs := []error(e)
	s := []string{}
	for _, e := range errs {
		s = append(s, e.Error())
	}

	return strings.Join(s, "\n")
}

// HasErr 是否有错误
func (e ResponseError) HasErr() bool {
	if len(e) == 0 {
		return false
	}
	return true
}

// Add 增加错误
func (e ResponseError) Add(err error) ResponseError {
	e = append(e, err)
	return e
}

// Response 回应对象
type Response struct {
	Request *Request
	Raw     *http.Response
	Body    []byte
	Errs    ResponseError
}

// Err 获取响应错误
func (r *Response) Err() error {
	if r.Errs.HasErr() {
		return r.Errs
	}
	return nil
}

// JSON 根据json绑定结构体
func (r *Response) JSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

// String 获取响应字符串
func (r *Response) String() string {
	return string(r.Body)
}

// Byte 获取响应字节
func (r *Response) Byte() []byte {
	return r.Body
}

// Status 获取响应状态码
func (r *Response) Status() int {
	return r.Raw.StatusCode
}

// Header 获取响应header
func (r *Response) Header() http.Header {
	return r.Raw.Header
}

// Cookies 获取响应 cookie
func (r *Response) Cookies() []*http.Cookie {
	return r.Raw.Cookies()
}

// IsError 是否响应错误
func (r *Response) IsError() bool {
	return r.Raw.StatusCode >= 400 && r.Err() != nil
}
