package tcpsrv

import (
  "fmt"
  "github.com/go-eyas/toolkit/log"
  "github.com/go-eyas/toolkit/tcp"
  "testing"
)

func TestClient(t *testing.T) {
  client, err := NewClientSrv(&tcp.Config{
    Addr:    ":6601",
    Logger: log.SugaredLogger,
  })
  if err != nil {
    panic(err)
  }

  client.On("register", func(response *TCPResponse) {
    fmt.Println("on receive register msg:", response)
  })

  client.On("userinfo", func(response *TCPResponse) {
    fmt.Println("on receive userinfo msg:", response)
  })

  res, err := client.Send("register", map[string]interface{}{
    "uid": 1234,
  })
  if err != nil {
    panic(err)
  }
  fmt.Println("send register response: ", res)

  res, err = client.Send("userinfo")
  if err != nil {
    panic(err)
  }
  fmt.Println("send userinfo response: ", res)

  // go func() {
  //   ch := client.Receive()
  //   for data := range ch {
  //     fmt.Println("receive data:", string(data.Data))
  //     // data.Response(map[string]interface{}{
  //     // 	"msg": "I'm fine.",
  //     // })
  //   }
  // }()

  // send register
  // data, _ := json.Marshal(map[string]interface{}{
  //   "cmd": "register",
  // })
  // err = client.Send(data)
  // if err != nil {
  //   panic(err)
  // }
  // time.Sleep(1 * time.Second)
  //
  // // send userinfo
  // data, _ = json.Marshal(map[string]interface{}{
  //   "cmd": "userinfo",
  // })
  // err = client.Send(data)
  // if err != nil {
  //   panic(err)
  // }
  //
  // c := make(chan bool, 0)
  // <- c
}
