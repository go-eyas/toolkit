package redis

import (
	"testing"
	"time"

	"github.com/go-eyas/toolkit/types"
)

func TestRedisConnError(t *testing.T) {
	_, err := New(&Config{})
	if err != nil {
		t.Log("redis conn error success")
	} else {
		panic("empty addrs should error")
	}
}
func TestRedis(t *testing.T) {
	r, err := New(&Config{
		Cluster: false,
		Addrs:   []string{"10.0.2.252:6379"},
		DB:      1,
		Prefix:  "test:prefix:",
	})
	if err != nil {
		panic("redis connect fail")
	}

	key := "tookit:test"
	val := `{"hello": "world"}`
	data := "world"

	// test Set
	err = r.Set(key, val)
	if err != nil {
		panic("set redis fail")
	}
	// test Get
	v, err := r.Get(key)
	if err != nil {
		panic("get redis fail")
	}

	// test Bind json
	res := struct {
		Hello string `json:"hello"`
	}{}
	err = types.JSONString(v).JSON(&res)

	if err == nil && res.Hello == data {
		t.Log("get redis success")
	} else {
		panic(err)
	}

	// test Del
	err = r.Del(key)
	if err != nil {
		panic("del redis key error")
	}
	v, err = r.Get(key)
	if err == nil && v == "" {
		t.Log("del key success")
	} else {
		panic("get redis fail")
	}

	// test sub/pub
	pbChan := make(chan *Message)
	go r.Sub("tookit-pub", func(msg *Message) {
		t.Logf("sub receive: %v", msg)
		pbChan <- msg
	})
	<-time.After(time.Second)
	err = r.Pub("tookit-pub", val)
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
