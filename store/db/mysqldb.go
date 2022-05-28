package db

import (
	"database/sql"
	"github.com/forwardalex/Ytool/assemblyInit"
	"github.com/forwardalex/Ytool/enum"
	"github.com/forwardalex/Ytool/layzeInit"
	"github.com/forwardalex/Ytool/model"
	"github.com/go-redis/redis/v8"
)

var conn *sql.DB

func GetConn() *sql.DB {
	if conn != nil {
		return conn
	}
	assembly := layzeInit.GetAssembly(enum.GetAssemblyEnum().MySQL)

	if assembly == nil {
		return nil
	}

	return assembly.(*sql.DB)
}

var (
	// RedisClient TODO
	RedisClient *redis.Client
)

func GetRedisConn() *redis.Client {
	if RedisClient != nil {
		return RedisClient
	}
	assembly := layzeInit.GetAssembly(enum.GetAssemblyEnum().Redis)

	if assembly == nil {
		return nil
	}

	return assembly.(*redis.Client)
}

func GetDBConfig() string {
	return assemblyInit.Dsn
}
func GetRedisConfig() model.RedisConf {
	return assemblyInit.GetRedisConfig()
}
