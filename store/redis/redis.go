package redis

import (
	"context"
	"fmt"
	"github.com/forwardalex/Ytool/log"
	"github.com/go-redis/redis/v8"
	"time"
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
	db:           0,       // use default db
	PoolSize:     100,
	MinIdleConns: 50,
}
*/
func NewClient(r *redis.Options) (err error) {
	client = redis.NewClient(r)
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}

//TryLock exp(ttl time)
func TryLock(ctx context.Context, key string, client *redis.Client, exp time.Duration) (*Lock, error) {
	locker := New(client)
	lock, err := locker.Obtain(ctx, key, exp, nil)
	if err == ErrNotObtained {
		fmt.Println("Could not obtain lock!")
		return nil, err
	} else if err != nil {
		log.Fatal(ctx, "", err.Error())
		return nil, err
	}
	return lock, nil
}
