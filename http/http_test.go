package http

import (
	"fmt"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	http := Header("Authorization", "Bearer asfdfadsfdsfdasfds").UseRequest(func(req *Request) *Request {
		fmt.Printf("http 发送 %s %s\n", req.SuperAgent.Method, req.SuperAgent.Url)
		return req
	}).UseResponse(func(req *Request, res *Response) *Response {
		fmt.Printf("http 接收 %s %s\n", req.SuperAgent.Method, req.SuperAgent.Url)
		return res
	}).Timeout(time.Second * 10).
		BaseURL("https://api.github.com").
		BaseURL("/repos")

	res, err := http.Get("/eyasliu/blog/issues", map[string]interface{}{
		"per_page": 1,
	})
	if err != nil {
		panic(err)
	}
	s := []struct {
		URL   string `json:"url"`
		Title string
	}{}
	err = res.JSON(&s)
	if err != nil {
		panic(err)
	}
	if len(s) > 0 {
		t.Logf("get res struct: %+v", s)
	}
}

func TestError(t *testing.T) {
	res, err := Get("https://api.github.com/repos/eyasliu/blog/issuesx", nil)
	if err != nil {
		t.Logf("get error success: statusCode=%d body=%s error=%s", res.Status(), res.String(), err.Error())
	} else {
		t.Fatalf("res: statusCode=%d body=%s", res.Status(), res.String())
		panic("should get 404 error")
	}

	res, err = Get("", nil)
	if err != nil {
		t.Logf("success empty url, statusCode=%d body=%s error=%s", res.Status(), res.String(), err.Error())
	} else {
		panic("should error")
	}

}

// func TestProxy(t *testing.T) {
// 	h := Proxy("http://127.0.0.1:1080")
// 	res, err := h.Get("https://www.google.com", map[string]string{
// 		"hl": "zh-Hans",
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	t.Logf("google html: %s", res.String())
// }
