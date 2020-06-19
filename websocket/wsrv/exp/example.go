package srv

import (
	"basic/config"
	"basic/util"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/go-eyas/toolkit/log"
	"github.com/go-eyas/toolkit/websocket"
)

type WSContext struct {
	MID          int64
	CMD          string
	Seqno        string
	RawData      json.RawMessage
	SessionID    uint64
	Socket       *websocket.Conn
	RawMessage   *websocket.Message
	Engine       *websocket.WS
	SrvClient    *Websocket
	Payload      []byte
	RequestData  *WSRequest
	ResponseData *WSResponse
}

func (c *WSContext) OK(args ...interface{}) {
	c.ResponseData.Status = util.CodeSuccess
	c.ResponseData.Msg = "ok"

	if len(args) > 0 {
		c.ResponseData.Data = args[0]
	}
}

func (c *WSContext) Bind(v interface{}) error {
	err := json.Unmarshal(c.RequestData.Data, v)
	if err != nil {
		return err
	}
	return validate.Struct(v)
}

var registerMu sync.Mutex

func (c *WSContext) Register(mid int64) {
	registerMu.Lock()
	c.SrvClient.session[c.SessionID] = mid
	registerMu.Unlock()
}

func (c *WSContext) writeResponse() error {
	payload, err := json.Marshal(c.ResponseData)
	if err != nil {
		return err
	}
	log.Debugf("[WS] --> 发送 MID=%d CMD=%s data=%s", c.MID, c.CMD, string(payload))
	return c.RawMessage.Response(payload)
}

type WSRequest struct {
	CMD   string          `json:"cmd"`
	Seqno string          `json:"seqno"`
	Data  json.RawMessage `json:"data"`
}

type WSResponse struct {
	CMD    string      `json:"cmd"`
	Seqno  string      `json:"seqno"`
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

type WSHandler func(*WSContext)

type Websocket struct {
	WS                  *websocket.WS
	routes              map[string]WSHandler
	session             map[uint64]int64 // map[sid]mid
	requestMiddlewares  []WSHandler      // 请求中间件
	responseMiddlewares []WSHandler      // 响应中间件
}

var WsSrv = &Websocket{}

// var WS *websocket.WS

func (ws *Websocket) Init(conf *config.Config) {
	WS := websocket.New(&websocket.Config{
		MsgType: websocket.BinaryMessage,
		// Logger:  log.SugaredLogger,
	})
	ws.WS = WS
	ws.routes = map[string]WSHandler{}
	ws.session = map[uint64]int64{}
	ws.requestMiddlewares = make([]WSHandler, 0)
	ws.responseMiddlewares = make([]WSHandler, 0)

	go ws.receive()
}

func (ws *Websocket) receive() {
	ch := ws.WS.Receive()
	for {
		req := <-ch
		go func(req *websocket.Message) {
			ctx := &WSContext{
				SessionID:   req.SID,
				Socket:      req.Socket,
				RawMessage:  req,
				Engine:      ws.WS,
				Payload:     req.Payload,
				RequestData: &WSRequest{},
				SrvClient:   ws,
			}
			ctx.MID = ws.session[req.SID]
			log.Debugf("[WS] <-- 收到 MID=%d CMD=%s data=%s", ctx.MID, ctx.CMD, string(ctx.Payload))
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
			}
			handler, ok := ws.routes[ctx.CMD]
			if !ok {
				return
			} else {
				ctx.ResponseData.Status = util.CodeEmptyImplement
				ctx.ResponseData.Msg = "empty implement"
			}
			handler(ctx)

			for _, mdl := range ws.responseMiddlewares {
				mdl(ctx)
			}
		}(req)

	}
}

func (ws *Websocket) Routes(routes map[string]WSHandler) {
	for k, h := range routes {
		fmt.Printf("[WS-debug] %-10s --> %s\n", k, util.FuncName(h))
		ws.routes[k] = h
	}
}

func (ws *Websocket) UseRequest(h WSHandler) {
	ws.requestMiddlewares = append(ws.requestMiddlewares, h)
}
func (ws *Websocket) UseResponse(h WSHandler) {
	ws.responseMiddlewares = append(ws.responseMiddlewares, h)
}

func (Websocket) Close() {}
