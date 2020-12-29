package http

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	http := New().TransformRequest(Logger.LoggerRequest).TransformResponse(Logger.LoggerResponse).Timeout(time.Second * 10).
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
	h := New().TransformRequest(Logger.LoggerRequest).TransformResponse(Logger.LoggerResponse)
	res, err := h.Header("just-test", "1234").
		Header("func-header", func() string {return "val in fn"}).
		Get("https://api.github.com/repos/eyasliu/blog/issuesx")
	if err == nil {
		panic("should get 404 error")
	}

	res, err = h.Get("", nil)
	if err != nil {
		t.Logf("success empty url, statusCode=%d body=%s error=%s", res.Status(), res.String(), err.Error())
	} else {
		panic("should error")
	}

	res, err = h.Post("https://api.github.com/repos/eyasliu/blog/issuesx", map[string]interface{}{"hello": "test"})
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

