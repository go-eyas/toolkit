package redis

import (
	"testing"
	"time"
)

func TestRedisConnError(t *testing.T) {
	err := Init(&Config{})
	if err != nil {
		t.Log("redis conn error success")
	} else {
		panic("empty addrs should error")
	}
}
func TestRedis(t *testing.T) {
	Init(&Config{
		Cluster: false,
		Addrs:   []string{"10.0.3.252:6379"},
		DB:      1,
	})

	key := "tookit:test"
	val := `{"hello": "world"}`
	data := "world"

	// test Set
	err := Set(key, val)
	if err != nil {
		panic("set redis fail")
	}
	// test Get
	v, err := Get(key)
	if err != nil {
		panic("get redis fail")
	}

	// test Bind json
	res := struct {
		Hello string `json:"hello"`
	}{}
	err = v.JSON(&res)

	if err == nil && res.Hello == data {
		t.Log("get redis success")
	} else {
		panic(err)
	}

	// test Del
	err = Del(key)
	if err != nil {
		panic("del redis key error")
	}
	v, err = Get(key)
	if err == nil && v.String() == "" {
		t.Log("del key success")
	} else {
		panic("get redis fail")
	}

	// test sub/pub
	pbChan := make(chan *Message)
	go Sub("tookit-pub", func(msg *Message) {
		t.Logf("sub receive: %v", msg)
		pbChan <- msg
	})
	<-time.After(time.Second)
	err = Pub("tookit-pub", val)
	if err != nil {
		panic(err)
	} else {
		t.Logf("pub success")
	}

	msg := <-pbChan

	res2 := struct {
		Hello string `json:"hello"`
	}{}
	err = msg.JSON(&res2)
	if err != nil || res2.Hello != data {
		panic("sub receive wrong")
	} else {
		t.Logf("sub receive success: %v", res2)
	}

}
