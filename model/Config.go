package model

//Conf 数据库配置
type DbConf struct {
	Host     string `json:"host"`     // 主机
	Port     int    `json:"port"`     // 端口
	User     string `json:"user"`     // 用户名
	Password string `json:"password"` // 用户密码
	Database string `json:"database"` // 数据库名
}

// Redis 配置
type RedisConf struct {
	Host     string `json:"host"`     // 主机
	Port     int    `json:"port"`     // 端口
	Password string `json:"password"` // 密码
	Type     string `json:"type"`     // 密码
}

type KafkaConf struct {
	KafkaServer string   `json:"kafka-server"` // kafka服务地址
	Topic       []string `json:"topic"`        // 主题
}

// grpc接口超时配置
type GrpcCallConf struct {
	ClientName string   `json:"clientName"`
	ServerName []string `json:"serverName"`
	Timeout    int      `json:"timeOut"`
}

type cosCfg struct {
	CDNHost   string `json:"cdn_host"`   // cdn加速访问地址
	AppID     string `json:"appid"`      // 腾讯云应用id
	Bucket    string `json:"bucket"`     // 存储桶名称，由用户自定义字符串和APPID组成
	Region    string `json:"region"`     // 地域，是腾讯云托管机房的分布地区
	SecretID  string `json:"secret_id"`  // 云API密钥，用于标识API调用者身份
	SecretKey string `json:"secret_key"` // 云API密钥，用于加密签名字符串和服务器端验证签名字符串的密钥
}

// cos配置
type CosConf struct {
	Public  *cosCfg `json:"public"`
	Private *cosCfg `json:"private"`
}

// 定时任务配置
type TaskConfig struct {
	Name string `json:"name"` // 定时任务名称
	Exec bool   `json:"exec"` // 是否执行
	Unit string `json:"unit"` // 时间单位
	Time uint64 `json:"time"` // 间隔时间
}
