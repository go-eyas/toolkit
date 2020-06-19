package wsrv

import (
	"encoding/json"
	"runtime/debug"

	"github.com/go-eyas/toolkit/log"
	"github.com/go-eyas/toolkit/util"
	"github.com/go-eyas/toolkit/websocket"
)

// WebsocketServer 服务器
type WebsocketServer struct {
	Engine              *websocket.WS
	logger              websocket.LoggerI
	routes              map[string][]WSHandler            // 路由
	session             map[uint64]map[string]interface{} // map[sid]sessionData
	requestMiddlewares  []WSHandler                       // 请求中间件
	responseMiddlewares []WSHandler                       // 响应中间件
}

// WSHandler 请求处理器
type WSHandler func(*Context)

// New 新建服务器实例
func New(conf *websocket.Config) *WebsocketServer {
	ws := websocket.New(conf)
	server := &WebsocketServer{
		Engine:              ws,
		routes:              make(map[string][]WSHandler),
		session:             make(map[uint64]map[string]interface{}),
		requestMiddlewares:  make([]WSHandler, 0),
		responseMiddlewares: make([]WSHandler, 0),
	}

	go server.receive()

	return server
}

func (ws *WebsocketServer) receive() {
	ch := ws.Engine.Receive()
	for {
		req := <-ch
		go func(req *websocket.Message) {
			ctx := &Context{
				SessionID:   req.SID,
				Socket:      req.Socket,
				RawMessage:  req,
				Engine:      ws.Engine,
				Payload:     req.Payload,
				RequestData: &WSRequest{},
				Server:      ws,
			}
			vals, ok := ws.session[req.SID]
			if !ok {
				ws.session[req.SID] = make(map[string]interface{})
				vals = ws.session[req.SID]
			}
			ctx.Values = vals
			log.Debugf("[WS] <-- 收到 CMD=%s data=%s", ctx.CMD, string(ctx.Payload))
			err := json.Unmarshal(ctx.Payload, ctx.RequestData)
			if err != nil {
				log.Errorf("WS request json parse error: %v", err)
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
					log.Errorf("%v", err)
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

func (ws *WebsocketServer) Handle(cmd string, handlers ...WSHandler) {
	h, ok := ws.routes[cmd]
	if !ok {
		h = handlers
	} else {
		h = append(h, handlers...)
	}
	ws.routes[cmd] = h
}
