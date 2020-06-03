package websocket

import (
	"net/http"
	"testing"
)

func TestWS(t *testing.T) {
	ws := New(&Config{
		Logger: logger,
	})

	http.HandleFunc("/ws", ws.HTTPHandler)
	http.HandleFunc("/", ws.Playground)

	go func() {
		rec := ws.Receive()
		for {
			req, _ := <-rec
			req.Response([]byte("1234556"))
		}
	}()

	// 浏览器打开 http://127.0.0.1:8800 测试
	t.Fatal(http.ListenAndServe("127.0.0.1:8800", nil))
}
