package tcpsrv

import (
  "fmt"
  "github.com/go-eyas/toolkit/log"
  "github.com/go-eyas/toolkit/tcp"
  "testing"
)

func TestServerSrv(t *testing.T) {
  srv, err := NewServerSrv(&tcp.Config{
    Addr:    ":6601",
    Logger: log.SugaredLogger,
  })

  if err != nil {
    panic(err)
  }

  srv.Use(func(c *Context) {
    fmt.Printf("TCP 收到 cmd=%s seqno=%s data=%s\n", c.CMD, c.Seqno, string(c.Payload))
    c.Next()
    fmt.Printf("TCP 响应 cmd=%s seqno=%s data=%s\n", c.CMD, c.Seqno, string(c.Response.Data))
  })
  srv.Use(func(c *Context) {
    if c.CMD != "register" {
      _, ok := c.Get("uid").(int64)
      if !ok {
        c.Response.Msg = "this connection is not register"
        c.Response.Status = 401
        c.Abort()
        return
      }
    }
    c.Next() // 如后续无操作，可省略
  })

  srv.Handle("register", func(c *Context) {
    body := &struct {
      UID int64 `json:"uid"`
    }{}
    err := c.Bind(body)
    if err != nil {
      panic(err)
    }
    c.Set("uid", body.UID)
    c.OK()
    c.Next()
  })

  srv.Handle("userinfo", func(c *Context) {
    uid := c.Get("uid").(int64)
    c.OK(uid)
  })

  c := make(chan bool, 0)
  <- c
}