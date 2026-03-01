package token_bucket

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/97wsn/gopkg/redis"

	"golang.org/x/time/rate"
)

func TestTokens_Allow(t1 *testing.T) {
	var limitTotal int32 = 0
	var luaTotal int32 = 0
	ctx := context.Background()
	rds := redis.NewTestRedis()
	tokens := NewTokens(rds)
	limiter := rate.NewLimiter(100, 100)
	g := sync.WaitGroup{}
	for n := 0; n < 10; n++ {
		if n%2 == 0 {
			time.Sleep(time.Duration(n*200) * time.Millisecond)
		}

		for i := 0; i < 3000; i++ {
			g.Add(2)
			go func(n1, n2 int) {
				defer g.Done()
				if tokens.Allow(ctx, &Request{
					Key:         "test-token",
					Burst:       100,
					LimitSecond: 100,
				}) {
					atomic.AddInt32(&luaTotal, 1)
				}
			}(n, i)

			go func(n1, n2 int) {
				defer g.Done()

				if limiter.Allow() {
					atomic.AddInt32(&limitTotal, 1)
				}
			}(n, i)
		}
	}

	g.Wait()

	fmt.Println(limitTotal, luaTotal)
}
