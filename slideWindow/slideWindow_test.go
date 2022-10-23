package slideWindow

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSlideWindow(t *testing.T) {
	swrl := NewSlidingWindowRateLimiter(&slidingWindowRateLimiterParam{
		slidingWindowLength:             5,
		timeSlice:                       time.Second,
		maxCountInTheWholeSlidingWindow: 10,
	})

	var wg sync.WaitGroup
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func(i int) {
			// defer recover防协程panic
			defer func() {
				e := recover()
				if e != nil {
					fmt.Errorf("panic here")
				}
			}()
			defer wg.Done()

			// 过一遍限流器的校验
			err := swrl.Take()
			if err != nil {
				fmt.Printf("slidingWindowRateLimiter Take error, num: %d\n", i)
				return
			}

			fmt.Println(i) // 执行业务逻辑的函数

			return
		}(i)
	}
	wg.Wait()
}
