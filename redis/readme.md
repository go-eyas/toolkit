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

```
package redis // import "github.com/go-eyas/toolkit/redis"


VARIABLES

var RedisTTL = time.Hour * 24
    RedisTTL 默认有效期 24 小时


FUNCTIONS

func Close()
    Close 关闭redis连接

func Del(keys ...string) error
    Del 删除键

func Expire(key string, expiration time.Duration) (bool, error)
func HDel(key string, field ...string) error
    HDel 删除hash的键

func HGetAll(key string) (map[string]StringValue, error)
    HGetAll 获取 Hash 的所有字段

func HSet(key, field string, val interface{}, expiration ...time.Duration) error
    HSet 设置hash值

func Init(redisConf *Config) error
    Init 初始化redis

func Pub(channel string, msg string) error
    Pub 发布事件 example: Redis.Pub("chat", "this is a test message")

func Set(key string, value interface{}, expiration ...time.Duration) error
    Set 设置字符串值，有效期默认 24 小时

func Sub(channel string, handler func(*Message))
    Sub 监听通道，有数据时触发回调 handler 
    example: 
    redis.Sub("chat")(func(msg *redis.Message) {
      fmt.Printf("receive message: %#v", msg)
    })

func (msg *Message) JSON(v interface{}) error
    JSON 绑定json对象

type RedisClient struct {
        Namespace string
        Client    redisClientInterface
        // Has unexported fields.
}
    RedisClient redis client wrapper

var Redis *RedisClient
    Redis 暴露的redis封装

type StringValue string
    StringValue redis 返回值

func Get(key string) (StringValue, error)
    Get 获取字符串值

func HGet(key string, field string) (StringValue, error)
    HGet 获取 Hash 的字段值

func (val StringValue) JSON(v interface{}) error
    JSON 将redis值转成指定结构体

func (val StringValue) String() string
    String 将redis值转成字符串


TYPES

type Config struct {
        Cluster  bool
        Addrs    []string
        Password string
        DB       int
}

type Message struct {
        Channel string
        Pattern string
        Payload StringValue
}

```