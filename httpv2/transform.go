package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type printLogger interface{
	Logf(string, ...interface{})
}

type consoleLogger struct {}
func (consoleLogger) Logf(s string, args ...interface{}) {
	fmt.Printf(s + "\n", args...)
}

type HttpLogger struct {
	logger printLogger
}

var Logger = &HttpLogger{logger: consoleLogger{}}

func (l *HttpLogger) LoggerRequest(r *Request, req *http.Request) *Request {
	var body []byte
	if req.Method != "GET" && req.Body != nil {
		body, _ = ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(io.Reader(bytes.NewReader(body)))
	}
	l.logger.Logf("HTTP SEND %s %s header=%v body=%s", req.Method, req.URL, req.Header, string(body))
	return r
}

func (l *HttpLogger) LoggerSafeResponse(r *Request, req *http.Request, resp *Response) *Response {
	body := "[body is unread]"
	if resp.IsRead {
		body = resp.String()
	}

	l.logger.Logf("HTTP RESV %s %s %d body=%s", req.Method, req.URL, resp.StatusCode(), body)
	return resp
}

func (l *HttpLogger) LoggerResponse(r *Request, req *http.Request, resp *Response) *Response {
	if resp.IsRead {
		resp.ReadAllBody()
	}
	l.logger.Logf("HTTP RESV %s %s %d %s", req.Method, req.URL, resp.StatusCode(), resp.String())
	return resp
}