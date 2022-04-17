package delayqueue

import (
	"context"
	"github.com/forwardalex/Ytool/log"
	"github.com/forwardalex/Ytool/store/db"
	"github.com/go-redis/redis/v8"
)

// BucketItem bucket中的元素
type BucketItem struct {
	timestamp int64
	jobId     string
}

// 添加JobId到bucket中
func pushToBucket(key string, timestamp int64, jobId string) error {
	return db.GetRedisConn().ZAdd(context.Background(), key, &redis.Z{Score: float64(timestamp),
		Member: jobId}).Err()
}

// 从bucket中获取延迟时间最小的JobId
func getFromBucket(key string) (*BucketItem, error) {
	cmd := db.GetRedisConn().ZRangeWithScores(context.Background(), key, 0, 0)
	value, err := cmd.Val(), cmd.Err()
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	if len(value) == 0 {
		return nil, nil
	}
	item := &BucketItem{}
	item.timestamp = int64(value[0].Score)
	item.jobId = value[0].Member.(string)
	ctx := context.Background()
	log.Info(ctx, "getFromBucket lasted item, timestamp=", item.timestamp, " jobId= ", item.jobId, " bucketname= ", key)
	return item, nil
}

// 从bucket中删除JobId
func removeFromBucket(bucket string, jobId string) error {
	return db.GetRedisConn().ZRem(context.Background(), bucket, jobId).Err()
}
