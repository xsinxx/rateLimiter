package leakyBucket

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLeakyBucket(t *testing.T) {
	RateLimiter := NewLeakyBucketRateLimiter(5, 1*time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func(i int) {
			defer func() {
				e := recover()
				if e != nil {
					fmt.Errorf("panic here")
				}
			}()
			defer wg.Done()

			RateLimiter.Take()
			fmt.Println(i) // 执行业务逻辑的函数
		}(i)
	}
	wg.Wait()
}
