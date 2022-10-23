package tokenBucket

import (
	"sync"
	"time"
)

type RateLimiter struct {
	rate   int64      // 令牌放入速度
	max    int64      // 令牌最大数量
	last   int64      // 上一次请求发生时间
	amount int64      // 令牌数量
	lock   sync.Mutex // 由于读写冲突，需要加锁
}

// 获得当前时间
func cur() int64 {
	return time.Now().Unix()
}

func NewTokenBucket(rate int64, max int64) *RateLimiter {
	return &RateLimiter{
		rate:   rate,
		max:    max,
		last:   cur(),
		amount: max,
	}
}

func (rl *RateLimiter) Take() bool {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	// 距离上一次请求过去的时间, 以秒为单位
	passed := cur() - rl.last

	// 计算在这段时间里 令牌数量可以增加多少
	amount := rl.amount + passed*rl.rate

	// 如果令牌数量超过上限；我们就不继续放入那么多令牌了
	if amount > rl.max {
		amount = rl.max
	}

	// 如果令牌数量仍然小于0，则说明请求应该拒绝
	if amount <= 0 {
		return false
	}

	amount--
	rl.amount = amount
	// 更新上次请求时间
	rl.last = cur()

	return true
}
