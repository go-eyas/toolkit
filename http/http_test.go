package http

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	http := Header("Authorization", "Bearer asfdfadsfdsfdasfds").UseRequest(func(req *Request) *Request {
    fmt.Printf("http 发送 %s %s\n", req.SuperAgent.Method, req.SuperAgent.Url)
    return req
	}).UseResponse(func(req *Request, res *Response) *Response {
	    fmt.Printf("http 接收 %s %s\n", req.SuperAgent.Method, req.SuperAgent.Url)
	    return res
	})

	res, err := http.Get("https://api.github.com/repos/eyasliu/blog/issues", map[string]interface{}{
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
