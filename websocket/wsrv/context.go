package wsrv

import (
  "encoding/json"
  "sync"

  "github.com/go-eyas/toolkit/gin/util"
  "github.com/go-eyas/toolkit/websocket"
  "github.com/go-playground/validator/v10"
)

var validate = validator.New()

// WSRequest 请求数据
type WSRequest struct {
  CMD   string          `json:"cmd"`
  Seqno string          `json:"seqno"`
  Data  json.RawMessage `json:"data"`
}

// WSResponse 响应数据
type WSResponse struct {
  CMD    string      `json:"cmd"`
  Seqno  string      `json:"seqno"`
  Status int         `json:"status"`
  Msg    string      `json:"msg"`
  Data   interface{} `json:"data"`
}

// Context 请求上下文
type Context struct {
  Values       map[string]interface{} // 该会话注册的值
  valMu        sync.RWMutex
  CMD          string             // 命令名称
  Seqno        string             // 请求唯一标识符
  RawData      json.RawMessage    // 请求原始数据 data
  SessionID    uint64             // 会话ID
  Socket       *websocket.Conn    // 长连接对象
  RawMessage   *websocket.Message // 原始消息对象
  Engine       *websocket.WS      // 引擎
  Server       *WebsocketServer   // 服务器对象
  Payload      []byte             // 请求原始消息报文
  RequestData  *WSRequest         // 已解析的请求数据
  ResponseData *WSResponse        // 响应数据
  logger       websocket.LoggerI
  isAbort      bool // 是否已停止继续执行中间件和处理函数
  sendMu sync.Mutex
}

// OK 响应成功数据
func (c *Context) OK(args ...interface{}) {
  c.ResponseData.Status = util.CodeSuccess
  c.ResponseData.Msg = "ok"

  if len(args) > 0 {
    c.ResponseData.Data = args[0]
  }
}

// Bind 解析并 JSON 绑定 data 数据到结构体，并验证数据正确性
func (c *Context) Bind(v interface{}) error {
  err := json.Unmarshal(c.RequestData.Data, v)
  if err != nil {
    return err
  }
  return validate.Struct(v)
}

// Get 获取会话的值
func (c *Context) Get(key string) interface{} {
  c.valMu.RLock()
  defer c.valMu.RUnlock()
  return c.Values[key]
}

// Set 设置会话的上下文的值，注意设置的值在整个会话生效，不仅仅在本次上下文请求而已
func (c *Context) Set(key string, v interface{}) {
  c.valMu.Lock()
  c.Values[key] = v
  c.valMu.Unlock()
}

// Abort 停止后面的处理函数和中间件执行
func (c *Context) Abort() {
  c.isAbort = true
}

// Push 服务器主动推送消息至该连接的客户端
func (c *Context) Push(data *WSResponse) error {
  return c.Server.Push(c.SessionID, data)
}

func (c *Context) writeResponse() error {
  c.sendMu.Lock()
  defer c.sendMu.Unlock()
  payload, err := json.Marshal(c.ResponseData)
  if err != nil {
    return err
  }
  c.logger.Infof("[WS] --> SEND CMD=%s data=%s", c.CMD, string(payload))
  return c.RawMessage.Response(payload)
}
