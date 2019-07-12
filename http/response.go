package http

import (
	"encoding/json"

	"github.com/parnurzeal/gorequest"
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
	s := ""
	for _, e := range errs {
		s += e.Error() + "\n"
	}

	return s
}

// HasErr 是否有错误
func (e ResponseError) HasErr() bool {
	if len(e) == 0 {
		return false
	}
	return true
}

// Response 回应对象
type Response struct {
	Request *Request
	Raw     *gorequest.Response
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
