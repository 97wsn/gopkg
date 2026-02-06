package token_bucket

import (
	"context"
	_ "embed"
	"time"

	"github.com/97wsn/gopkg/redis"
)

type Request struct {
	// 请求的key
	Key string `json:"key"`
	// 允许的容量
	Burst int `json:"burst"`
	// 每秒发送多少个
	LimitSecond float64 `json:"limit_second"`
	// 当前消耗token数
	UseNum int `json:"use_num"`
}

type Result struct {
	// 存在的tokens数
	Tokens float64 `db:"tokens"`
	// 最后处理时间(纳秒)
	LastTime int64 `db:"last_time"`
}

//go:embed script.lua
var luaScript string

type Tokens struct {
	cache *redis.Redis
}

func NewTokens(rds *redis.Redis) *Tokens {
	return &Tokens{cache: rds}
}

// Allow 判断是否允许请求
// 具体实现可以参考: golang.org/x/time/rate
// See https://en.wikipedia.org/wiki/Token_bucket for more about token buckets.
func (t *Tokens) Allow(ctx context.Context, request *Request) bool {
	result, err := t.cache.Client().Eval(ctx,
		luaScript,
		[]string{t.cache.PrefixKey(request.Key)},
		request.Burst, request.LimitSecond, max(request.UseNum, 1), time.Now().UnixNano(),
	).Result()
	if err != nil {
		// 注意: Redis出现问题不能限制流量
		return true
	}

	return result == int64(1)
}
