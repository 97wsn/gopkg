package redis

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func NewTestRedis() *Redis {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	return NewRedis(redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	}))
}
