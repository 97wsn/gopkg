package fixed_window

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gopkg/redis"
)

const redisPrefix = "rate-limiter:"

type Limiter struct {
	rds *redis.Redis
}

func NewLimiter(rds *redis.Redis) *Limiter {
	return &Limiter{rds: rds}
}

type options struct {
	num      int
	interval time.Duration
}

type Option func(*options)

func WithNum(n int) Option {
	return func(o *options) {
		o.num = n
	}
}

func WithInterval(t time.Duration) Option {
	return func(o *options) {
		o.interval = t
	}
}

func (l *Limiter) Allow(ctx context.Context, key string, opts ...Option) bool {
	opt := &options{
		num:      1000,
		interval: time.Second * 10,
	}

	for _, o := range opts {
		o(opt)
	}

	key = fmt.Sprintf("%s:%s", redisPrefix, key)
	if lt := l.rds.TTL(ctx, key); lt <= 0 {
		err := l.rds.SetEX(ctx, key, strconv.Itoa(opt.num), opt.interval)
		if err != nil {
			return true
		}
	}

	num, err := l.rds.Decr(ctx, key)
	if err != nil {
		return true
	}

	return num > 0
}

// Reset gets a key and reset all limitations and previous usages
func (l *Limiter) Reset(ctx context.Context, key string) error {
	return l.rds.Del(ctx, redisPrefix+key)
}
