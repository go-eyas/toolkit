package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ResponseError struct {
	errs []error
}
func (e *ResponseError) Error() string {
	msgs := ""
	for _, err := range e.errs {
		if len(msgs) > 0 {
			msgs += "\n"
		}
		msgs += err.Error()
	}
	return msgs
}

func (e *ResponseError) Add(err error) {
	if err != nil {
		e.errs = append(e.errs, err)
	}
}

type Response struct {
	Client *Client
	Request *http.Request
	Response *http.Response
	body []byte
	Err *ResponseError
	IsRead bool
}

func newResponse(request *Client, r *http.Request) *Response {
	return &Response{
		Client: request,
		Request: r,
		Err: &ResponseError{errs: make([]error, 0)},
	}
}

func (rp *Response) ready() {
	//if rp.Response == nil {
	//	return
	//}
	if code := rp.StatusCode(); code >= 400 {
		rp.AddError(fmt.Errorf("http status code %d", code))
	}
}

// StatusCode 获取 HTTP 响应状态码
func (rp *Response) StatusCode() int {
	if rp.Response == nil {
		return 0
	}
	return rp.Response.StatusCode
}

// Status StatusCode 别名
func (rp *Response) Status() int {
	return rp.StatusCode()
}

// GetError 获取错误
func (rp *Response) GetError() error {
	if len(rp.Err.errs) > 0 {
		return rp.Err
	}
	return nil
}

// AddError 手动增加错误
func (rp *Response) AddError(err error) {
	rp.Err.Add(err)
}

// ReadAllBody 读取响应的 Body 流
func (rp *Response) ReadAllBody() (bt []byte, err error) {
	if rp.Response != nil && !rp.IsRead {
		bt, err = ioutil.ReadAll(rp.Response.Body)
		rp.body = bt
		rp.IsRead = true
		return
	}
	bt = rp.body
	return
}

// Body 获取响应的原始数据
func (rp *Response) Body() (bt []byte) {
	bt, _ = rp.ReadAllBody()
	return
}

// SetBody 重置响应的 Body
func (rp *Response) SetBody(bt []byte) {
	rp.IsRead = true
	rp.body = bt
}

// String 将响应的 Body 数据转成字符串
func (rp *Response) String() string {
	if !rp.IsRead { rp.ReadAllBody() }
	return string(rp.Body())
}

// Error 实现 error interface
func (rp *Response) Error() string {
	return rp.Err.Error()
}

// JSON 使用 JSON 解析响应 Body 数据
func (rp *Response) JSON(v interface{}) error {
	if !rp.IsRead { rp.ReadAllBody() }
	return json.Unmarshal(rp.body, v)
}