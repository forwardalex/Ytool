package debug

import (
	"context"
	"fmt"
	"github.com/forwardalex/Ytool/log"
	"github.com/forwardalex/Ytool/model"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type debugConfig struct {
	Debug struct {
		ENV       string `yaml:"env"`
		NeedDebug bool   `yaml:"needDebug"`
		Mail      struct {
			User     string `yaml:"user"`
			Host     string `yaml:"host"`
			Password string `yaml:"password"`
			Port     int    `yaml:"port"`
		} `yaml:"mail"`
		Etcd struct {
			Endpoints            []string      `yaml:"endpoints"`
			AutoSyncInterval     time.Duration `yaml:"auto-sync-interval"`
			DialTimeout          int           `yaml:"dialTimeout "`
			DialKeepAliveTime    time.Duration `yaml:"dial-keep-alive-time"`
			DialKeepAliveTimeout time.Duration `yaml:"dial-keep-alive-timeout"`
			MaxCallSendMsgSize   int
			MaxCallRecvMsgSize   int
		} `yaml:"etcd"`
		Mysql struct {
			Host     string `yaml:"host"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Port     int    `yaml:"port"`
			Database string `yaml:"database"`
		} `yaml:"mysql"`
		Redis struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Password string `yaml:"password"`
			Type     string `yaml:"type"`
		} `yaml:"redis"`
		Mock        []mockRpcConfig `yaml:"mock"`
		MockContext struct {
			Head struct {
				Seq     int    `yaml:"seq"`
				Uin     int64  `yaml:"uin"`
				OpenId  string `yaml:"openid"`
				UserId  string `yaml:"userid"`
				TraceId string `yaml:"traceid"`
				ExtData string `yaml:"extdata"`
				Client  struct {
					IP   string `yaml:"ip"`
					Port int    `yaml:"port"`
				} `yaml:"client"`
			} `yaml:"head"`
		} `yaml:"mockcontext"`

		MockConfig []mockConfig `yaml:"mockPrivateConfig"`
	}
}

type mockRpcConfig struct {
	Clientname string `yaml:"clientname"`
	Servername string `yaml:"servername"`
	Response   string `yaml:"response"`
}

type mockConfig struct {
	Key  string `yaml:"key"`
	Conf string `yaml:"conf"`
}

var Config debugConfig
var NeedDebug bool

var filePath = "./work-config.yml"

func init() {
	log.Init("dev")
	fmt.Println("-----检测当前环境-----")
	// 如果是开发环境，且需要开启调试，则读取配置文件
	//先设置dev 读取配置yaml文件配置
	os.Setenv("ENV_NAME", "Dev")
	getenv := os.Getenv("ENV_NAME")
	getenv = strings.ToUpper(getenv)
	if strings.ToUpper(getenv) == "DEV" {
		exists, err := PathExists(filePath)
		if err != nil || !exists {
			fmt.Println("找不到配置文件，不执行本地debug模式")
			return
		}

		// 加载配置文件
		Config.GetConf()
		os.Setenv("ENV_NAME", Config.Debug.ENV)
		log.Info(context.Background(), "ENV_NAME", os.Getenv("ENV_NAME"))
		if Config.Debug.NeedDebug {
			fmt.Println("------本地环境可以开启debug模式------")
			NeedDebug = true
		}
	}
}

// 在开发环境进行DB重定向
func ConfigDebugDB(dbConf *model.DbConf) {
	if NeedDebug {

		if &Config.Debug == nil || &Config.Debug.Mysql == nil || Config.Debug.Mysql.Host == "" {
			fmt.Println("-----配置文件没有数据库链接重定向配置，使用七彩石配置-----")
			return
		}

		fmt.Println("-----开始进行数据库连接重定向-----")
		dbDebug := Config.Debug.Mysql
		dbConf.Host = dbDebug.Host
		dbConf.User = dbDebug.User
		dbConf.Password = dbDebug.Password
		dbConf.Port = dbDebug.Port
		dbConf.Database = dbDebug.Database
		fmt.Println("-----数据库连接重定向成功-----")
	}
}

// 在开发环境对Redis重定向
func ConfigDebugRedis(redisConf *model.RedisConf) {
	if NeedDebug {

		if &Config.Debug == nil || &Config.Debug.Redis == nil || Config.Debug.Redis.Host == "" {
			fmt.Println("-----配置文件没有Redis连接重定向配置，使用七彩石配置-----")
			return
		}

		fmt.Println("-----开始进行Redis连接重定向-----")
		redisDebug := Config.Debug.Redis
		redisConf.Host = redisDebug.Host
		redisConf.Port = redisDebug.Port
		redisConf.Password = redisDebug.Password
		redisConf.Type = redisDebug.Type
		fmt.Println("-----Redis连接重定向成功-----")
	}
}

func (c *debugConfig) GetConf() *debugConfig {
	err := ReadYml(filePath, &c)
	if err != nil {
		fmt.Println(err)
	}
	return c
}

func PathExists(path string) (bool, error) {

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//读取配置文件
func ReadYml(path string, obj interface{}) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return yaml.Unmarshal(yamlFile, obj)
}

func GetMailConf() (mail *model.MailConf) {
	mail = &model.MailConf{
		User:     Config.Debug.Mail.User,
		PassWord: Config.Debug.Mail.Password,
		Host:     Config.Debug.Mail.Host,
		Port:     Config.Debug.Mail.Port,
	}
	return mail
}
func GetEtcdConf() (etcd *model.EtcdConf) {
	etcd = &model.EtcdConf{
		Endpoints:   Config.Debug.Etcd.Endpoints,
		DialTimeout: time.Duration(Config.Debug.Etcd.DialTimeout) * time.Second,
	}
	return etcd
}
