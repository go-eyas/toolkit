package tcpsrv

import (
  "encoding/json"
  "github.com/go-eyas/toolkit/gin/util"
  "github.com/go-eyas/toolkit/tcp"
  "github.com/go-playground/validator/v10"
  "reflect"
  "sync"
)

var validate = validator.New()

// WSRequest 请求数据
type TCPRequest struct {
  CMD   string          `json:"cmd"`
  Seqno string          `json:"seqno"`
  Data  json.RawMessage `json:"data"`
}

func (r *TCPRequest) SetJSON(v interface{}) (err error) {
  r.Data, err = convertDataToJsonByte(v)
  return
}

func (r *TCPRequest) BindJSON(v interface{}) error {
  err := json.Unmarshal(r.Data, v)
  if err != nil {
    return err
  }
  rt := reflect.TypeOf(v)
  if rt.Kind() == reflect.Ptr {
    rt = rt.Elem()
  }
  if rt.Kind() == reflect.Struct {
    return validate.Struct(v)
  }
  return nil
}

// WSResponse 响应数据
type TCPResponse struct {
  CMD    string          `json:"cmd"`
  Seqno  string          `json:"seqno"`
  Status int             `json:"status"`
  Msg    string          `json:"msg"`
  Data   json.RawMessage `json:"data"`
}

func (r *TCPResponse) SetJSON(v interface{}) (err error) {
  r.Data, err = convertDataToJsonByte(v)
  return
}

func (r *TCPResponse) BindJSON(v interface{}) error {
  err := json.Unmarshal(r.Data, v)
  if err != nil {
    return err
  }
  rt := reflect.TypeOf(v)
  if rt.Kind() == reflect.Ptr {
    rt = rt.Elem()
  }
  if rt.Kind() == reflect.Struct {
    return validate.Struct(v)
  }
  return nil
}

// Context 请求上下文
type Context struct {
  Values       map[string]interface{} // 该会话注册的值
  valMu        sync.RWMutex
  CMD          string          // 命令名称
  Seqno        string          // 请求唯一标识符
  RawData      json.RawMessage // 请求原始数据 data
  SessionID    uint64          // 会话ID
  Socket       *tcp.Conn       // 长连接对象
  RawMessage   *tcp.Message    // 原始消息对象
  Engine       *tcp.Server     // 引擎
  Server       *ServerSrv      // 服务器对象
  Payload      []byte          // 请求原始消息报文
  Request  *TCPRequest     // 已解析的请求数据
  Response *TCPResponse    // 响应数据
  logger   tcp.LoggerI
  handlers []TCPHandler // 当前请求上下文的处理器
  handlerIndex int          // 当前中间件处理
  isAbort      bool         // 是否已停止继续执行中间件和处理函数
  sendMu       sync.Mutex
}

// OK 响应成功数据
func (c *Context) OK(args ...interface{}) error {
  c.Response.Status = util.CodeSuccess
  c.Response.Msg = "ok"

  if len(args) > 0 {
    return c.Response.SetJSON(args[0])
  } else {
    c.Response.Data = []byte("null")
    return nil
  }

}

// Bind 解析并 JSON 绑定 data 数据到结构体，并验证数据正确性
func (c *Context) Bind(v interface{}) error {
  return c.Request.BindJSON(v)
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
func (c *Context) Push(data *TCPResponse) error {
  return c.Server.Push(c.SessionID, data)
}

func (c *Context) Next() {
  if c.isAbort {
    return
  }
  if c.handlerIndex < len(c.handlers)-1 {
    c.handlerIndex++
    c.handlers[c.handlerIndex](c)
  } else {
    c.Abort()
  }
}

func (c *Context) writeResponse() error {
  c.sendMu.Lock()
  defer c.sendMu.Unlock()
  payload, err := json.Marshal(c.Response)
  if err != nil {
    return err
  }
  c.logger.Infof("[TCP] --> SEND CMD=%s data=%s", c.CMD, string(payload))
  return c.RawMessage.Response(payload)
}

func convertDataToJsonByte(data interface{}) ([]byte, error) {
  var bodyData []byte
  if data == nil {
    bodyData = []byte("null")
  } else if _bodyData, ok := data.([]byte); ok {
    bodyData = _bodyData
  } else if _bodyData, ok := data.(string); ok {
    bodyData = []byte(_bodyData)
  } else {
    _bodyData, err := json.Marshal(data)
    if err != nil {
      return nil, err
    }
    bodyData = _bodyData
  }
  return bodyData, nil
}