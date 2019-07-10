package redis

import (
	"testing"
)

func TestRedisConnError(t *testing.T) {
	err := Init(&RedisConfig{})
	if err != nil {
		t.Log("redis conn error success")
	} else {
		panic("empty addrs should error")
	}
}
func TestRedis(t *testing.T) {
	Init(&RedisConfig{
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
}
