package wsrv

import (
  "encoding/json"
  "fmt"
  "runtime/debug"
  "sync"
  "time"

  "github.com/go-eyas/toolkit/util"
  "github.com/go-eyas/toolkit/websocket"
)

// WebsocketServer 服务器
type WebsocketServer struct {
  Engine              *websocket.WS
  Config              *websocket.Config
  logger              websocket.LoggerI
  routes              map[string][]WSHandler            // 路由
  Session             map[uint64]map[string]interface{} // map[sid]SessionData
  heartbeat           *sync.Map            // 心跳
  requestMiddlewares  []WSHandler                       // 请求中间件
  responseMiddlewares []WSHandler                       // 响应中间件
}

// WSHandler 请求处理器
type WSHandler func(*Context)

var sessionMu sync.Mutex

// New 新建服务器实例
func New(conf *websocket.Config) *WebsocketServer {
  if conf.Logger == nil {
    conf.Logger = websocket.EmptyLogger
  }
  if conf.MsgType == 0 {
    conf.MsgType = websocket.BinaryMessage
  }
  ws := websocket.New(conf)
  server := &WebsocketServer{
    Config:              conf,
    logger:              conf.Logger,
    Engine:              ws,
    routes:              make(map[string][]WSHandler),
    Session:             make(map[uint64]map[string]interface{}),
    heartbeat:           &sync.Map{},
    requestMiddlewares:  make([]WSHandler, 0),
    responseMiddlewares: make([]WSHandler, 0),
  }

  ws.HandleClose(server.onClose)
  ws.HandleCreate(server.onCreate)

  go server.receive()
  go server.checkHeartbeat()

  return server
}

func (ws *WebsocketServer) receive() {
  ch := ws.Engine.Receive()
  for {
    req := <-ch
    go func(req *websocket.Message) {
      ws.heartbeat.Store(req.SID, time.Now().Unix())

      // 心跳包
      if len(req.Payload) == 0 {
        return
      }

      ctx := &Context{
        SessionID:   req.SID,
        Socket:      req.Socket,
        RawMessage:  req,
        Engine:      ws.Engine,
        Payload:     req.Payload,
        RequestData: &WSRequest{},
        Server:      ws,
        logger:      ws.logger,
      }
      sessionMu.Lock()
      vals, ok := ws.Session[req.SID]
      if !ok {
        ws.Session[req.SID] = make(map[string]interface{})
        vals = ws.Session[req.SID]
      }
      sessionMu.Unlock()
      ctx.Values = vals
      ctx.logger.Infof("[WS] <-- RECV CMD=%s data=%s", ctx.CMD, string(ctx.Payload))

      err := json.Unmarshal(ctx.Payload, ctx.RequestData)
      if err != nil {
        ctx.logger.Errorf("WS request json parse error: %v", err)
        return
      }
      ctx.CMD = ctx.RequestData.CMD
      ctx.Seqno = ctx.RequestData.Seqno
      ctx.ResponseData = &WSResponse{
        CMD:    ctx.RequestData.CMD,
        Seqno:  ctx.RequestData.Seqno,
        Status: -1,
        Msg:    "not implement",
        Data:   map[string]interface{}{},
      }

      defer func() {
        if err := recover(); err != nil {
          ws.logger.Errorf("%v", err)
          debug.PrintStack()
          r := util.ParseError(err)
          ctx.ResponseData.Status = r.Status
          ctx.ResponseData.Msg = r.Msg
          ctx.ResponseData.Data = r.Data
        }
        ctx.writeResponse()
      }()

      for _, mdl := range ws.requestMiddlewares {
        mdl(ctx)
        if ctx.isAbort {
          break
        }
      }
      handler, ok := ws.routes[ctx.CMD]
      if !ok {
        return
      } else {
        ctx.ResponseData.Status = 1
        ctx.ResponseData.Msg = "empty implement"
      }
      if !ctx.isAbort {
        for _, h := range handler {
          h(ctx)
          if ctx.isAbort {
            break
          }
        }
      }

      for _, mdl := range ws.responseMiddlewares {
        mdl(ctx)
        if ctx.isAbort {
          break
        }
      }
    }(req)
  }
}

// UseRequest 请求中间件
func (ws *WebsocketServer) UseRequest(h WSHandler) {
  ws.requestMiddlewares = append(ws.requestMiddlewares, h)
}

// UseResponse 响应中间件
func (ws *WebsocketServer) UseResponse(h WSHandler) {
  ws.responseMiddlewares = append(ws.responseMiddlewares, h)
}

// Handle 注册 CMD 路由监听器
func (ws *WebsocketServer) Handle(cmd string, handlers ...WSHandler) {
  h, ok := ws.routes[cmd]
  if !ok {
    h = handlers
  } else {
    h = append(h, handlers...)
  }
  ws.routes[cmd] = h
}

// Push 服务器推送消息到客户端
func (ws *WebsocketServer) Push(sid uint64, data *WSResponse) error {
  conn, ok := ws.Engine.Clients[sid]
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
  return conn.Send(&websocket.Message{
    SID:     sid,
    Payload: payload,
    Socket:  conn,
    MsgType: ws.Config.MsgType,
  })
}

// Destroy 销毁清理连接
func (ws *WebsocketServer) Destroy(sid uint64) {
  conn, ok := ws.Engine.Clients[sid]
  if ok {
    conn.Destroy()
  }
  ws.heartbeat.Delete(conn.ID)
  sessionMu.Lock()
  delete(ws.Session, sid)
  sessionMu.Unlock()
}

func (ws *WebsocketServer) onCreate(conn *websocket.Conn) {
  sid := conn.ID
  ws.heartbeat.Store(sid, time.Now().Unix())
  sessionMu.Lock()
  ws.Session[sid] = make(map[string]interface{})
  sessionMu.Unlock()
}

func (ws *WebsocketServer) onClose(conn *websocket.Conn) {
  sid := conn.ID
  ws.Destroy(sid)
}

func (ws *WebsocketServer) checkHeartbeat() {
  for {
    time.Sleep(time.Second)
    now := time.Now().Unix() - 30
    ws.heartbeat.Range(func(key, val interface{}) bool {
      sid, ok := key.(uint64)
      if !ok {
        return true
      }
      hbTime, ok := val.(int64)
      if !ok {
        return true
      }

      if hbTime < now {
        go ws.Destroy(sid)
      }
      return true
    })
  }
}