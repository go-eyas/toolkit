# redis


## 初始化
只支持单例 redis

```go
import "github.com/go-eyas/toolkit/redis"

func main() {
  // 使用前必须先初始化
  err := redis.Init(&redis.Config{
    Cluster:  false, // 是否集群
    Addrs:    []string{"10.0.3.252:6379"}, // redis 地址，如果是集群则在数组上写多个元素
    Password: "",
		DB:       1,
  })
  if err != nil {
    panic(err)
  }

  err = redis.Set("tookit:test", `{"hello": "world"}`)

  v, err = redis.Get("tookit:test")
  v.String() // {"hello": "world"}

  data := struct {
		Hello string `json:"hello"`
	}{}
  v.JSON(&data)
  data.Hello // "world"

  err = redis.Del("tookit:test")


  redis.Expire("tookit:test", time.Hour * 24)

  redis.Redis // *redis.RedisClient
  redis.Client // *github.com/go-redis/redis.Client
}
```

## godoc

[API 文档](https://gowalker.org/github.com/go-eyas/toolkit/redis)