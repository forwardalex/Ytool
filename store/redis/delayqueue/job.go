package delayqueue

import (
	"context"
	"encoding/json"
	"github.com/forwardalex/Ytool/store/db"
)

// Job 使用msgpack序列化后保存到Redis,减少内存占用
type Job struct {
	Topic string `json:"topic"`
	Id    string `json:"id"`    // job唯一标识ID
	Delay int64  `json:"delay"` // 延迟时间, unix时间戳
	TTR   int64  `json:"ttr"`   // Job执行超时时间, 单位：秒
	Body  string `json:"body"`
}

// 获取Job
func getJob(key string) (*Job, error) {
	cmd := db.GetRedisConn().Get(context.Background(), key)
	value, err := cmd.Val(), cmd.Err()
	if err != nil {
		return nil, err
	}
	if len(value) == 0 {
		return nil, nil
	}

	JsonStr := []byte(value)
	job := &Job{}
	err = json.Unmarshal(JsonStr, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

// 添加Job
func putJob(key string, job Job) error {
	value, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return db.GetRedisConn().Set(context.Background(), key, string(value), 0).Err()
}

// 删除Job
func removeJob(key string) error {
	err := db.GetRedisConn().Del(context.Background(), key).Err()
	return err
}
