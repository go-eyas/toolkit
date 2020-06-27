package tcpsrv

import (
  "encoding/json"
  "errors"
  "github.com/go-eyas/toolkit/emit"
  "github.com/go-eyas/toolkit/tcp"
  "github.com/go-eyas/toolkit/util"
  "sync"
  "time"
)

// WSRequest 请求数据
type TCPClientRequest struct {
  CMD   string      `json:"cmd"`
  Seqno string      `json:"seqno"`
  Data  interface{} `json:"data"`
}

type ClientSrv struct {
  sendTimeout time.Duration
  Engine      *tcp.Client
  Emitter     *emit.Emitter
  sendMutex   sync.Mutex
  sendProcess map[string]chan *TCPResponse
}

// NewClientSrv 实例化客户端服务
func NewClientSrv(conf *tcp.Config) (*ClientSrv, error) {
  engine, err := tcp.NewClient(conf)
  if err != nil {
    return nil, err
  }
  srv := &ClientSrv{
    Engine:      engine,
    Emitter:     emit.New(),
    sendTimeout: 10 * time.Second,
    sendProcess: make(map[string]chan *TCPResponse),
  }

  go srv.reader()

  return srv, nil
}

func (cs *ClientSrv) reader() {
  ch := cs.Engine.Receive()
  for msg := range ch {
    res := &TCPResponse{}
    err := json.Unmarshal(msg.Data, res)
    if err != nil {
      cs.Emitter.Emit("error", err)
      continue
    }
    cs.sendMutex.Lock()
    ch, ok := cs.sendProcess[res.Seqno]
    cs.sendMutex.Unlock()
    if ok {
      ch <- res
    }
    cs.Emitter.Emit(res.CMD, res)
  }
}

// On 监听服务器响应数据，每当服务器有数据发送过来，都会以 cmd 为事件名触发监听函数
func (cs *ClientSrv) On(cmd string, h func(*TCPResponse)) {
  cs.Emitter.On(cmd, func(_res interface{}) {
    res, ok := _res.(*TCPResponse)
    if ok {
      h(res)
    }
  })
}

// Pub 给服务器发送消息
func (cs *ClientSrv) Pub(cmd string, data interface{}) error {
  _, err := cs.writeSend(cmd, data)
  return err
}

// Send 给服务器发送消息，并等待服务器的响应数据，10秒超时
func (cs *ClientSrv) Send(cmd string, datas ...interface{}) (*TCPResponse, error) {
  var data interface{}
  if len(datas) > 0 {
    data = datas[0]
  }
  body, err := cs.writeSend(cmd, data)
  if err != nil {
    return nil, err
  }
  cs.sendMutex.Lock()
  cs.sendProcess[body.Seqno] = make(chan *TCPResponse)
  cs.sendMutex.Unlock()
  ticker := time.Tick(cs.sendTimeout)
  select {
  case res := <-cs.sendProcess[body.Seqno]:
    cs.sendMutex.Lock()
    close(cs.sendProcess[body.Seqno])
    delete(cs.sendProcess, body.Seqno)
    cs.sendMutex.Unlock()
    return res, nil
  case <-ticker:
    cs.sendMutex.Lock()
    close(cs.sendProcess[body.Seqno])
    delete(cs.sendProcess, body.Seqno)
    cs.sendMutex.Unlock()
    return nil, errors.New("request timeout")
  }

}

// 封装数据，并发送数据到服务端
func (cs *ClientSrv) writeSend(cmd string, data interface{}) (*TCPRequest, error) {
  body := &TCPRequest{
    CMD:   cmd,
    Seqno: util.RandomStr(8),
  }
  err := body.SetJSON(data)
  if err != nil {
    return nil, err
  }

  raw, err := json.Marshal(body)
  if err != nil {
    return nil, err
  }
  return body, cs.Engine.Send(raw)
}
