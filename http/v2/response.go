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

func (rp *Response) StatusCode() int {
	if rp.Response == nil {
		return 0
	}
	return rp.Response.StatusCode
}

func (rp *Response) Status() int {
	return rp.StatusCode()
}

func (rp *Response) GetError() error {
	if len(rp.Err.errs) > 0 {
		return rp.Err
	}
	return nil
}

func (rp *Response) AddError(err error) {
	rp.Err.Add(err)
}

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

func (rp *Response) Body() (bt []byte) {
	bt, _ = rp.ReadAllBody()
	return
}

func (rp *Response) String() string {
	if !rp.IsRead { rp.ReadAllBody() }
	return string(rp.Body())
}

func (rp *Response) Error() string {
	return rp.Err.Error()
}

func (rp *Response) JSON(v interface{}) error {
	if !rp.IsRead { rp.ReadAllBody() }
	return json.Unmarshal(rp.body, v)
}