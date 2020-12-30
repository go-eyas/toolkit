package http_test

import (
	"github.com/go-eyas/toolkit/log"
	"github.com/go-eyas/toolkit/http/v2"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	h := http.Use(http.AccessLogger(log.SugaredLogger)).Timeout(time.Second * 10).
		BaseURL("https://api.github.com").
		BaseURL("/repos")

	res, err := h.Get("/eyasliu/blog/issues", map[string]interface{}{
		"per_page": 1,
	})
	if err != nil {
		panic(err)
	}
	var s []struct {
		URL   string `json:"url"`
		Title string
	}
	err = res.JSON(&s)
	if err != nil {
		panic(err)
	}
	if len(s) > 0 {
		t.Logf("get res struct: %+v", s)
	}
}

func TestError(t *testing.T) {
	h := http.New().Use(http.AccessLogger(log.SugaredLogger))
	res, err := h.Header("just-test", "1234").
		Header("func-header", func() string {return "val in fn"}).
		Get("https://api.github.com/repos/eyasliu/blog/issuesx")
	if err == nil {
		panic("should get 404 error")
	}

	res, err = h.Header("a", "1").Get("", nil)
	if err != nil {
		t.Logf("success empty url, statusCode=%d body=%s error=%s", res.Status(), res.String(), err.Error())
	} else {
		panic("should error")
	}

	res, err = h.Header("b", "2").Type("form").Post("https://api.github.com/repos/eyasliu/blog/issuesx", map[string]interface{}{"hello": "test"})
	if err != nil {
		t.Logf("success empty url, statusCode=%d body=%s error=%s", res.Status(), res.String(), err.Error())
	} else {
		panic("should error")
	}
	h.Safe(false)
	h.Header("c", "3")
	res, err = h.Header("d", "4").Post("https://api.github.com/repos/eyasliu/blog/issuesx", map[string]interface{}{"hello": "test"})
	if err != nil {
		t.Logf("success empty url, statusCode=%d body=%s error=%s", res.Status(), res.String(), err.Error())
	} else {
		panic("should error")
	}
	h = h.Safe(true).BaseURL("http://notexistdomain.qwer")
	h.Type("form").Post("/", `file=@file:./http.go`)

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

