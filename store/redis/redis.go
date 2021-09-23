package redis

import (
	"github.com/go-redis/redis"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

type (
	Redis struct {
		Addr string
		Type string
		Pass string
		tls  bool
	}
	Option func(r *redis.Options)
)

// NewClient 推荐链接格式
/* 	&redis.Options{
	Addr:         "127.0.0.1:6379",
	Password:     "hello", // no password set
	DB:           0,       // use default DB
	PoolSize:     100,
	MinIdleConns: 50,
}
*/
func NewClient(r *redis.Options) (err error) {
	client = redis.NewClient(r)
	_, err = client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
