package assemblyInit

import (
	"context"
	"fmt"
	"github.com/forwardalex/Ytool/debug"
	"github.com/forwardalex/Ytool/enum"
	"github.com/forwardalex/Ytool/log"
	"github.com/forwardalex/Ytool/model"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql" //mysql 驱动
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type RedisInit struct {
	obj model.AssemblyObj
}

func (impl *RedisInit) InitAssembly(ctx context.Context) interface{} {
	db, err := initRedis()
	if err != nil {
		log.Fatal(context.Background(), "", err.Error())
	}
	return db
}

func (*RedisInit) GetAssemblyType() enum.Enum {
	return enum.GetAssemblyEnum().Redis
}

func (impl *RedisInit) GetAssemblyObj() *model.AssemblyObj {
	return &impl.obj
}

//InitDB 初始化数据库
func initRedis() (*redis.Client, error) {
	var (
		redisbConf model.RedisConf
		err        error
	)
	ctx := context.Background()
	err = getRedisConfig(&redisbConf)

	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	fmt.Println(path[:index])

	// 如果是开发环境，手动指定数据库地址
	debug.ConfigDebugRedis(&redisbConf)
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     redisbConf.Host + ":" + strconv.Itoa(redisbConf.Port),
		Password: redisbConf.Password,
	})
	if _, err = RedisClient.Ping(ctx).Result(); err != nil {
		log.Error(ctx, "redis connect error:", err.Error())
		return nil, err
	}
	log.Info(ctx, "redis connect ok")
	return RedisClient, nil
}

func getRedisConfig(conf *model.RedisConf) error {
	return nil
}
