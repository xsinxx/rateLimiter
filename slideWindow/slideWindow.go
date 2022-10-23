package slideWindow

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type slidingWindowRateLimiterParam struct {
	slidingWindowLength             int           // 有几个窗口
	timeSlice                       time.Duration // 每个窗口时间多长
	maxCountInTheWholeSlidingWindow int           // 全部窗口整体能承受的总流量上限
}

type slidingWindowRateLimiter struct {
	param slidingWindowRateLimiterParam // 参数

	countChan  chan int64 // 窗口队列，维护当前各窗口区间计数
	curCount   int64      // 当前窗口计数。
	totalCount int64      // 全部计数
	onceLock   sync.Once  // 保证Tick函数只跑一次
}

func NewSlidingWindowRateLimiter(param *slidingWindowRateLimiterParam) *slidingWindowRateLimiter {
	return &slidingWindowRateLimiter{
		param:      *param,
		countChan:  make(chan int64, param.slidingWindowLength-1),
		curCount:   0,
		totalCount: 0,
	}
}

// Take 增加计数, 达到限额则返回错误
func (limiter *slidingWindowRateLimiter) Take() error {
	limiter.onceLock.Do(func() {
		go limiter.Tick()
	})
	if atomic.LoadInt64(&limiter.totalCount) >= int64(limiter.param.maxCountInTheWholeSlidingWindow) {
		return errors.New("reach limit")
	}
	atomic.AddInt64(&limiter.curCount, 1)
	atomic.AddInt64(&limiter.totalCount, 1)
	return nil
}

// Tick 窗口滑动, 总计数更改
func (limiter *slidingWindowRateLimiter) Tick() {
	for {
		select {
		case <-time.After(limiter.param.timeSlice):
			if len(limiter.countChan) == cap(limiter.countChan) {
				c := <-limiter.countChan
				atomic.AddInt64(&limiter.totalCount, -c)
			}
			if len(limiter.countChan) < cap(limiter.countChan) {
				c := atomic.SwapInt64(&limiter.curCount, 0) // 返回的limiter.curCount中的值
				limiter.countChan <- c
			}
		}
	}
}
