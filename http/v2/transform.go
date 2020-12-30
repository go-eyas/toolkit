package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type printLogger interface{
	Debug(...interface{})
}

type consoleLogger struct {}
func (consoleLogger) Debug(a ...interface{}) {
	fmt.Println(a...)
}

// 打印 HTTP 请求响应日志
type HttpLogger struct {
	Logger           printLogger
	MaxLoggerBodyLen int64
}

func AccessLogger(logger printLogger, bodyLimit ...int64) *HttpLogger {
	var limit int64 = 2048
	if len(bodyLimit) > 0 {
		limit = bodyLimit[0]
	}
	return &HttpLogger{Logger: logger, MaxLoggerBodyLen: limit}
}

var Logger = AccessLogger(consoleLogger{})

// 打印 HTTP 请求日志
func (l *HttpLogger) TransformRequest(c *Client, req *http.Request) *Client {
	var body []byte
	logtext := fmt.Sprintf("HTTP SEND %s %s header=%v", req.Method, req.URL, req.Header)
	if req.Method != "GET" && req.Body != nil  {
		logtext = fmt.Sprintf("%s size=%d", logtext, req.ContentLength)
		// 如果body太长，估计是文件上传，不打印，也不侵入，并且太长的body也会妨碍控制台输出
		if req.ContentLength < l.MaxLoggerBodyLen {
			body, _ = ioutil.ReadAll(req.Body)
			req.Body = ioutil.NopCloser(io.Reader(bytes.NewReader(body)))
			logtext = fmt.Sprintf("%s body=%s", logtext, string(body))
		}
	}
	l.Logger.Debug(logtext)
	return c
}

// 打印 HTTP 响应日志
// warning: 会先读取一遍 Response.Body，如果该中间件导致了 http 下载异常问题，请关闭该中间件
func (l *HttpLogger) TransformResponse(c *Client, req *http.Request, resp *Response) *Response {
	logText := fmt.Sprintf("HTTP RESV %s %s %d", req.Method, req.URL, resp.StatusCode())
	if resp.IsRead {
		resp.ReadAllBody()
	}
	logText += " " + resp.String()
	l.Logger.Debug(logText)
	return resp
}
