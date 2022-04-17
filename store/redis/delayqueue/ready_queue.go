package delayqueue

import (
	"context"
	"fmt"
	"github.com/forwardalex/Ytool/store/db"
	"time"
)

// 添加JobId到队列中
func pushToReadyQueue(queueName string, jobId string) error {
	queueName = fmt.Sprintf(Setting.QueueName, queueName)
	err := db.GetRedisConn().RPush(context.Background(), queueName, jobId).Err()

	return err
}

// 从队列中阻塞获取JobId
func blockPopFromReadyQueue(queues []string, timeout int) (string, error) {
	var args []string
	for _, queue := range queues {
		queue = fmt.Sprintf(Setting.QueueName, queue)
		args = append(args, queue)
	}

	cmd := db.GetRedisConn().BLPop(context.Background(), time.Duration(timeout*int(time.Second)), args...)
	value, err := cmd.Val(), cmd.Err()
	if err != nil {
		return "", err
	}
	if value == nil {
		return "", nil
	}
	if len(value) == 0 {
		return "", nil
	}

	return value[1], nil
}
