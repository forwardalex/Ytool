package model

import "time"

type EtcdConf struct {
	Endpoints            []string      `json:"endpoints"`
	AutoSyncInterval     time.Duration `json:"auto-sync-interval"`
	DialTimeout          time.Duration `json:"dial-timeout"`
	DialKeepAliveTime    time.Duration `json:"dial-keep-alive-time"`
	DialKeepAliveTimeout time.Duration `json:"dial-keep-alive-timeout"`
	MaxCallSendMsgSize   int
	MaxCallRecvMsgSize   int
}
