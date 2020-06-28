package tcpsrv

import (
  "encoding/json"
  "fmt"
  "github.com/go-eyas/toolkit/tcp"
  "github.com/go-eyas/toolkit/util"
  "runtime/debug"
  "sync"
  "time"
)

type TCPHandler func(*Context)

type ServerSrv struct {
  Engine              *tcp.Server
  Config *tcp.Config
  logger tcp.LoggerI
  routes map[string][]TCPHandler           // 路由
  Session             map[uint64]map[string]interface{} // map[sid]SessionData
  sessionMu           sync.Mutex
  heartbeat           *sync.Map    // 心跳
  handlerMiddlewares []TCPHandler // 中间件
}

func NewServerSrv(conf *tcp.Config) (*ServerSrv, error) {
  engine, err := tcp.NewServer(conf)
  if err != nil {
    return nil, err
  }
  srv := &ServerSrv{
    Engine:              engine,
    Config:              conf,
    logger:              engine.Config.Logger,
    routes:              make(map[string][]TCPHandler),
    Session:             make(map[uint64]map[string]interface{}),
    heartbeat:           &sync.Map{},
    handlerMiddlewares:  make([]TCPHandler, 0),
  }

  srv.Engine.HandleCreate(srv.onCreate)
  srv.Engine.HandleClose(srv.onClose)

  go srv.receive()
  go srv.checkHeartbeat()

  return srv, nil
}

func (srv *ServerSrv) onCreate(conn *tcp.Conn) {
  sid := conn.ID
  srv.heartbeat.Store(sid, time.Now().Unix())
  srv.sessionMu.Lock()
  srv.Session[sid] = make(map[string]interface{})
  srv.sessionMu.Unlock()
  srv.logger.Infof("[TCP] New conn id=%d", conn.ID)
}
func (srv *ServerSrv) onClose(conn *tcp.Conn) {
  sid := conn.ID
  srv.Destroy(sid)
  srv.logger.Infof("[TCP] CLOSE conn id=%d", conn.ID)
}

func (srv *ServerSrv) handlerReceive(req *tcp.Message) {
  conn := req.Conn
  srv.heartbeat.Store(conn.ID, time.Now().Unix())

  // 心跳包
  if len(req.Data) == 0 {
    return
  }


  ctx := &Context{
    SessionID:   conn.ID,
    Socket:      conn,
    RawMessage:  req,
    Engine:      srv.Engine,
    Payload:     req.Data,
    Request: &TCPRequest{},
    Server:      srv,
    logger:      srv.logger,
  }
  srv.sessionMu.Lock()
  vals, ok := srv.Session[conn.ID]
  if !ok {
    srv.Session[conn.ID] = make(map[string]interface{})
    vals = srv.Session[conn.ID]
  }
  srv.sessionMu.Unlock()
  ctx.Values = vals
  ctx.logger.Infof("[TCP] <-- RECV CMD=%s data=%s", ctx.CMD, string(ctx.Payload))

  err := json.Unmarshal(ctx.Payload, ctx.Request)
  if err != nil {
    ctx.logger.Errorf("TCP request json parse error: %v", err)
    return
  }
  ctx.CMD = ctx.Request.CMD
  ctx.Seqno = ctx.Request.Seqno
  ctx.Response = &TCPResponse{
    CMD:    ctx.Request.CMD,
    Seqno:  ctx.Request.Seqno,
    Status: -1,
    Msg:    "not implement",
    // Data:   map[string]interface{}{},
  }

  defer func() {
      if err := recover(); err != nil {
        srv.logger.Errorf("%v", err)
        debug.PrintStack()
        r := util.ParseError(err)
        ctx.Response.Status = r.Status
        ctx.Response.Msg = r.Msg
        ctx.Response.SetJSON(r.Data)
      }
      ctx.writeResponse()
  }()

  handlers := append([]TCPHandler{}, srv.handlerMiddlewares...)
  handler, ok := srv.routes[ctx.CMD]
  if ok {
    handlers = append(handlers, handler...)
  }
  ctx.handlers = handlers
  ctx.handlerIndex = -1

  for !ctx.isAbort && ctx.handlerIndex < len(ctx.handlers) {
    ctx.Next()
  }
}

func (srv *ServerSrv) receive()        {
  ch := srv.Engine.Receive()
  for {
    req := <-ch
    go srv.handlerReceive(req)
  }
}

// 处理器中间件
func (srv *ServerSrv) Use(h ...TCPHandler) {
  srv.handlerMiddlewares = append(srv.handlerMiddlewares, h...)
}

// Handle 注册 CMD 路由监听器
func (srv *ServerSrv) Handle(cmd string, handlers ...TCPHandler) {
  h, ok := srv.routes[cmd]
  if !ok {
    h = handlers
  } else {
    h = append(h, handlers...)
  }
  srv.routes[cmd] = h
}

// Push 服务器推送消息到客户端
func (srv *ServerSrv) Push(sid uint64, data *TCPResponse) error {
  conn, ok := srv.Engine.Sockets[sid]
  if !ok {
    return fmt.Errorf("sid=%d is invalid", sid)
  }
  if data.Seqno == "" {
    data.Seqno = util.RandomStr(8)
  }
  if data.Status == 0 && data.Msg == "" {
    data.Msg = "ok"
  }

  payload, err := json.Marshal(data)
  if err != nil {
    return err
  }
  return conn.Send(payload)
}

func (srv *ServerSrv) checkHeartbeat() {
  for {
    time.Sleep(time.Second)
    now := time.Now().Unix() - 30
    srv.heartbeat.Range(func(key, val interface{}) bool {
      sid, ok := key.(uint64)
      if !ok {
        return true
      }
      hbTime, ok := val.(int64)
      if !ok {
        return true
      }

      if hbTime < now {
        go srv.onClose(srv.Engine.Sockets[sid])
      }
      return true
    })
  }
}
func (srv *ServerSrv) Destroy(sid uint64) {
  conn, ok := srv.Engine.Sockets[sid]
  if !ok {
    return
  }
  srv.heartbeat.Delete(conn.ID)
  srv.sessionMu.Lock()
  delete(srv.Session, sid)
  srv.sessionMu.Unlock()
}
