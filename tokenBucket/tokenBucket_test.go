package tokenBucket

import (
	"fmt"
	"sync"
	"testing"
)

func TestTokenBucket(t *testing.T) {
	RateLimiter := NewTokenBucket(3, 10)

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

			status := RateLimiter.Take()
			if !status {
				fmt.Printf("slidingWindowRateLimiter Take error, num: %d\n", i)
				return
			}
			fmt.Println(i) // 执行业务逻辑的函数
		}(i)
	}
	wg.Wait()
}
