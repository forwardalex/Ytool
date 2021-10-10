package env

import (
	"errors"
	"fmt"
	"github.com/forwardalex/Ytool/log"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

// Env 定义环境变量用词
type Env string

var Getenv string

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

func Service() string {
	return os.Getenv("service")
}

//HostIP 公网ip
func HostIP() string {
	//call ifconfig 查询ip
	// local ip
	//	conn, err := net.Dial("udp", "www.baidu.com:80")
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		return
	//	}
	//	defer conn.Close()
	//  return conn.LocalAddr().String()

	// 公网ip
	responseClient, errClient := http.Get("https://ipw.cn/api/ip/myip") // 获取外网 IP
	if errClient != nil {
		log.Error("err ", "获取外网 IP 失败，请检查网络")
		panic(errClient)
	}
	// 程序在使用完 response 后必须关闭 response 的主体。
	defer responseClient.Body.Close()

	body, _ := ioutil.ReadAll(responseClient.Body)
	clientIP := fmt.Sprintf("%s", string(body))
	return clientIP
}

func PodName() string {
	return os.Getenv("podName")
}

func GetENV() string {
	if Getenv != "" {
		return Getenv
	} else {
		Getenv = os.Getenv("ENV_NAME")
		return Getenv
	}
}

//InternalIp 内网ip
func InternalIp() string {
	infs, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, inf := range infs {
		if isEthDown(inf.Flags) || isLoopback(inf.Flags) {
			continue
		}

		addrs, err := inf.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}

	return ""
}
func isEthDown(f net.Flags) bool {
	return f&net.FlagUp != net.FlagUp
}

func isLoopback(f net.Flags) bool {
	return f&net.FlagLoopback == net.FlagLoopback
}
