package leakyBucket

import (
	"sync"
	"time"
)

type Bucket struct {
	bucketChan chan interface{} // 漏桶找总容量
	rate       time.Duration    // 速率
	onceLock   sync.Once        // 保证tick函数仅执行一次
}

func NewLeakyBucketRateLimiter(maxSize int64, rate time.Duration) *Bucket {
	return &Bucket{
		bucketChan: make(chan interface{}, maxSize),
		rate:       rate,
	}
}

func (b *Bucket) Take() {
	b.onceLock.Do(func() {
		go b.Tick()
	})
	<-b.bucketChan // 有缓冲channel取不到数据则等待
}

func (b *Bucket) Tick() {
	for {
		select {
		case <-time.After(b.rate):
			if len(b.bucketChan) != cap(b.bucketChan) {
				b.bucketChan <- 1
			}
		}
	}
}
