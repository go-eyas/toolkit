package wsrv

import (
	"net/http"
	"testing"

	"github.com/go-eyas/toolkit/log"
	"github.com/go-eyas/toolkit/websocket"
)

func TestSrv(t *testing.T) {
	server := New(&websocket.Config{
		Logger: log.SugaredLogger,
	})
	server.UseRequest(func(c *Context) {
		log.Debugf("ws request middleware, sid=%d", c.SessionID)
	})
	server.UseResponse(func(c *Context) {
		log.Debugf("ws response middleware, sid=%d", c.SessionID)
	})
	server.UseRequest(func(c *Context) {
		uid, ok := c.Get("uid").(int64)
		if !ok || uid == 0 {
			c.Abort()
		}
	})
	server.Handle("register")
	server.Handle("register", func(c *Context) {
		c.Set("uid", int(123))
		c.OK()
	})

	http.HandleFunc("/ws", server.Engine.HTTPHandler)
	http.HandleFunc("/play", server.Engine.Playground)
	http.ListenAndServe("127.0.0.1:9000", nil)

}
