package delayqueue

import (
	"errors"
	"fmt"
	"log"
	"time"
)

const (
	// DefaultBucketSize bucket数量
	DefaultBucketSize = 3
	// DefaultBucketName bucket名称
	DefaultBucketName = "dq_bucket_%d"
	// DefaultQueueName 队列名称
	DefaultQueueName = "dq_queue_%s"
	// DefaultQueueBlockTimeout 轮询队列超时时间
	DefaultQueueBlockTimeout = 180
	// DefaultBucketTicker 扫描bucket定时器时间间隔
	DefaultBucketTicker = 5
)

var (
	// 每个定时器对应一个bucket
	timers []*time.Ticker
	// bucket名称chan
	bucketNameChan <-chan string
	// Setting
	// 延迟队列配置项
	Setting *Config
)

// Config 应用配置
type Config struct {
	BucketSize        int    // bucket数量
	BucketName        string // bucket在redis中的键名,
	QueueName         string // ready queue在redis中的键名
	QueueBlockTimeout int    // 调用blpop阻塞超时时间, 单位秒, 修改此项, redis.read_timeout必须做相应调整
	BucketTicker      int    // 扫描bucket定时器时间间隔,单位秒
}

// Start  开启延时队列
// Init 初始化延时队列
func Start(cfg *Config) {
	initDefaultConfig(cfg)
	initTimers()
	bucketNameChan = generateBucketName()
}

// 初始化默认配置
func initDefaultConfig(cfg *Config) {
	Setting = &Config{}
	if cfg == nil {
		Setting.BucketSize = DefaultBucketSize
		Setting.BucketName = DefaultBucketName
		Setting.QueueName = DefaultQueueName
		Setting.QueueBlockTimeout = DefaultQueueBlockTimeout
		Setting.BucketTicker = DefaultBucketTicker
		return
	}

	if cfg.BucketSize <= 0 || cfg.BucketSize > 10000 {
		Setting.BucketSize = DefaultBucketSize
	}

	if len(cfg.BucketName) == 0 {
		Setting.BucketName = DefaultBucketName
	}

	if len(cfg.QueueName) == 0 {
		Setting.QueueName = DefaultQueueName
	}

	if cfg.QueueBlockTimeout == 0 {
		Setting.QueueBlockTimeout = DefaultQueueBlockTimeout
	}
	if cfg.BucketTicker == 0 {
		Setting.BucketTicker = DefaultBucketTicker
	}
}

// Push 添加一个Job到队列中
func Push(job Job) error {
	if job.Id == "" || job.Topic == "" || job.Delay < 0 || job.TTR <= 0 {
		return errors.New("invalid job")
	}

	err := putJob(job.Id, job)
	if err != nil {
		log.Printf("添加job到job pool失败#job-%+v#%s", job, err.Error())
		return err
	}
	err = pushToBucket(<-bucketNameChan, job.Delay, job.Id)
	if err != nil {
		log.Printf("添加job到bucket失败#job-%+v#%s", job, err.Error())
		return err
	}

	return nil
}

// Pop 轮询获取Job
func Pop(topics []string) (*Job, error) {
	jobId, err := blockPopFromReadyQueue(topics, Setting.QueueBlockTimeout)
	if err != nil {
		return nil, err
	}

	// 队列为空
	if jobId == "" {
		return nil, nil
	}

	// 获取job元信息
	job, err := getJob(jobId)
	if err != nil {
		return job, err
	}

	// 消息不存在, 可能已被删除
	if job == nil {
		return nil, nil
	}

	// 任务执行完成后需调用finish接口删除任务, 否则任务会重复投递, 消费端需能处理同一任务的多次投递
	// timestamp := time.Now().Unix() + job.TTR
	// err = pushToBucket(<-bucketNameChan, timestamp, job.Id)

	return job, err
}

// Remove 删除Job
func Remove(jobId string) error {
	return removeJob(jobId)
}

// Get 查询Job
func Get(jobId string) (*Job, error) {
	job, err := getJob(jobId)
	if err != nil {
		return job, err
	}

	// 消息不存在, 可能已被删除
	if job == nil {
		return nil, nil
	}
	return job, err
}

// 初始化定时器
func initTimers() {
	timers = make([]*time.Ticker, Setting.BucketSize)
	var bucketName string
	for i := 0; i < Setting.BucketSize; i++ {
		timers[i] = time.NewTicker(time.Duration(Setting.BucketTicker) * time.Second)
		bucketName = fmt.Sprintf(Setting.BucketName, i+1)
		go waitTicker(timers[i], bucketName)
	}
}

func waitTicker(timer *time.Ticker, bucketName string) {
	for {
		select {
		case t := <-timer.C:
			tickHandler(t, bucketName)
		}
	}
}

//todo  这里轮询不适合多实例  改分布式
// 扫描bucket, 取出延迟时间小于当前时间的Job
func tickHandler(t time.Time, bucketName string) {
	for {
		bucketItem, err := getFromBucket(bucketName)
		if err != nil {
			log.Printf("扫描bucket错误#bucket-%s#%s", bucketName, err.Error())
			return
		}

		// 集合为空
		if bucketItem == nil {
			return
		}

		// 延迟时间未到
		if bucketItem.timestamp > t.Unix() {
			return
		}

		// 延迟时间小于等于当前时间, 取出Job元信息并放入ready queue
		job, err := getJob(bucketItem.jobId)
		if err != nil {
			log.Printf("获取Job元信息失败#bucket-%s#%s", bucketName, err.Error())
			continue
		}

		// job元信息不存在, 从bucket中删除
		if job == nil {
			removeFromBucket(bucketName, bucketItem.jobId)
			continue
		}

		// 再次确认元信息中delay是否小于等于当前时间
		if job.Delay > t.Unix() {
			// 从bucket中删除旧的jobId
			removeFromBucket(bucketName, bucketItem.jobId)
			// 重新计算delay时间并放入bucket中
			pushToBucket(<-bucketNameChan, job.Delay, bucketItem.jobId)
			continue
		}

		err = pushToReadyQueue(job.Topic, bucketItem.jobId)
		if err != nil {
			log.Printf("JobId放入ready queue失败#bucket-%s#job-%+v#%s",
				bucketName, job, err.Error())
			continue
		}

		// 从bucket中删除
		removeFromBucket(bucketName, bucketItem.jobId)
	}
}

// 轮询获取bucket名称, 使job分布到不同bucket中, 提高扫描速度
func generateBucketName() <-chan string {
	c := make(chan string)
	go func() {
		i := 1
		for {
			c <- fmt.Sprintf(Setting.BucketName, i)
			if i >= Setting.BucketSize {
				i = 1
			} else {
				i++
			}
		}
	}()

	return c
}
