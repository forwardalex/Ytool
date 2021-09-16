package configs

import (
	"errors"
	"fmt"
	"os"
)

// Env 定义环境变量用词
type Env string

const (
	// EnvLocal 本地调试环境
	EnvLocal Env = "local"
	// EnvDevelopment 开发环境
	EnvDevelopment Env = "dev"
	// EnvTest 测试环境
	EnvTest Env = "test"
	// EnvPreRelease 预发布环境
	EnvPreRelease Env = "pre"
	// EnvProduction 生产环境
	EnvProduction Env = "prod"
)

// EnvMapStr 环境映射字符串
var EnvMapStr = map[Env]string{
	EnvLocal:       "local",
	EnvDevelopment: "dev",
	EnvTest:        "test",
	EnvPreRelease:  "pre",
	EnvProduction:  "prod",
}

// StrMapEnv 字符串映射环境
var StrMapEnv = map[string]Env{
	"local": EnvLocal,
	"dev":   EnvDevelopment,
	"test":  EnvTest,
	"pre":   EnvPreRelease,
	"prod":  EnvProduction,
}

// EnvMapSuffixStr 环境后缀
var EnvMapSuffixStr = map[Env]string{
	EnvLocal:       "", // 本地调试环境不需要后缀
	EnvDevelopment: "dev",
	EnvTest:        "test",
	EnvPreRelease:  "pre",
	EnvProduction:  "", // 正式环境不需要后缀
}

// Current 获取当前环境
func Current() Env {
	var (
		env Env
		ok  bool
	)
	svrEnv := os.Getenv("SVR_ENV")
	if svrEnv == "" {
		panic(errors.New("SVR_ENV is empty"))
	}
	if env, ok = StrMapEnv[svrEnv]; !ok {
		panic(fmt.Errorf("unknown env: %s", svrEnv))
	}
	return env
}
