package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Server       string
	Password     string
	DB           int
	Prefix       string `json:"prefix" toml:"prefix" yaml:"prefix"` // 统一前缀
	MaxRetries   int    `json:"max_retries" toml:"max_retries" yaml:"max_retries"`
	OpenTracing  bool   `json:"opentracing" toml:"opentracing" yaml:"opentracing"` // 是否开启链路追踪
	DialTimeout  int    `json:"dial_timeout" toml:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  int    `json:"read_timeout" toml:"read_timeout" yaml:"read_timeout"`
	WriteTimeout int    `json:"write_timeout" toml:"write_timeout" yaml:"write_timeout"`
}

// Open redis client
func Open(addr string, options ...func(options *redis.Options)) *redis.Client {
	redisOption := &redis.Options{
		Addr: addr,
	}

	for _, option := range options {
		option(redisOption)
	}

	return redis.NewClient(redisOption)
}

func Connect(conf *Config) (*redis.Client, error) {
	r := newRedis(conf)
	log.Println("[redis] connect:" + conf.Server)

	_, err := r.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	client := newRedis(conf)

	return client, nil
}

// 创建 redis for config
func newRedis(conf *Config) *redis.Client {
	return Open(conf.Server, func(options *redis.Options) {
		options.Password = conf.Password
		options.DB = conf.DB

		if conf.MaxRetries > 0 {
			options.MaxRetries = conf.MaxRetries
		}

		if conf.DialTimeout > 0 {
			options.DialTimeout = time.Duration(conf.DialTimeout) * time.Second
		}

		if conf.ReadTimeout > 0 {
			options.ReadTimeout = time.Duration(conf.ReadTimeout) * time.Second
		}

		if conf.WriteTimeout > 0 {
			options.WriteTimeout = time.Duration(conf.WriteTimeout) * time.Second
		}
	})
}
