package fixed_window

import (
	"context"
	"testing"
	"time"

	"gopkg/redis"
)

func TestLimiter_Allow(t *testing.T) {
	rds := redis.NewTestRedis()
	limiter := NewLimiter(rds)

	for i := 0; i < 9; i++ {
		t.Log(i)
		if !limiter.Allow(context.Background(), "test", WithNum(10), WithInterval(10*time.Second)) {
			t.Fatal("Allow failed 1")
		}
	}

	if limiter.Allow(context.Background(), "test", WithNum(10), WithInterval(10*time.Second)) {
		t.Fatal("Allow failed 2")
	}
}
