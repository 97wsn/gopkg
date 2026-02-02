package redis

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	CACHE_WEEK_TTL = 604800 * time.Second // 7天
	CACHE_DAY_TTL  = 86400 * time.Second  // 1天
	CACHE_HOUR_TTL = 3600 * time.Second   // 1小时
	CACHE_5MIN_TTL = 300 * time.Second    // 5分钟
	CACHE_MIN_TTL  = 60 * time.Second     // 1分钟
)

type PubSub = redis.PubSub

type Redis struct {
	client *redis.Client
	prefix string
}

type Option func(*Redis)

func WithPrefix(prefix string) Option {
	return func(r *Redis) {
		r.prefix = prefix
	}
}

func NewRedis(client *redis.Client, opts ...Option) *Redis {
	r := &Redis{
		client: client,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (b *Redis) Client() *redis.Client {
	return b.client
}

func (b *Redis) PrefixKey(key string) string {
	return b.prefix + key
}

func (b *Redis) Del(ctx context.Context, key string) error {
	return b.client.Del(ctx, b.PrefixKey(key)).Err()
}

// Set Our specification is that all keys must have an expiration time, so the Set key will expire in one day by default
func (b *Redis) Set(ctx context.Context, key, val string) error {
	return b.client.Set(ctx, b.PrefixKey(key), val, CACHE_HOUR_TTL).Err()
}

func (b *Redis) Unlink(ctx context.Context, key string) error {
	return b.client.Unlink(ctx, b.PrefixKey(key)).Err()
}

func (b *Redis) SetEX(ctx context.Context, key, val string, expiration time.Duration) error {
	return b.client.Set(ctx, b.PrefixKey(key), val, expiration).Err()
}

func (b *Redis) SetNX(ctx context.Context, key, val string, expiration time.Duration) (bool, error) {
	return b.client.SetNX(ctx, b.PrefixKey(key), val, expiration).Result()
}

func (b *Redis) Get(ctx context.Context, key string) (string, error) {
	return b.client.Get(ctx, b.PrefixKey(key)).Result()
}

func (b *Redis) MGet(ctx context.Context, key ...string) (map[string]string, error) {
	keys := make([]string, 0, len(key))

	for _, v := range key {
		keys = append(keys, b.PrefixKey(v))
	}

	result, err := b.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	results := make(map[string]string, len(key))
	for i, v := range result {
		if v == nil {
			continue
		}
		results[key[i]] = v.(string)
	}
	return results, nil
}

func (b *Redis) Incr(ctx context.Context, key string) (int64, error) {
	return b.client.Incr(ctx, b.PrefixKey(key)).Result()
}

func (b *Redis) Decr(ctx context.Context, key string) (int64, error) {
	return b.client.Decr(ctx, b.PrefixKey(key)).Result()
}

func (b *Redis) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return b.client.IncrBy(ctx, b.PrefixKey(key), value).Result()
}

func (b *Redis) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return b.client.DecrBy(ctx, b.PrefixKey(key), value).Result()
}

func (b *Redis) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return b.client.Expire(ctx, b.PrefixKey(key), expiration).Result()
}

func (b *Redis) HGetAll(ctx context.Context, key string, m any) error {
	info, err := b.HGetAllMap(ctx, key)
	if err != nil {
		return err
	}
	return scanStruct(info, m)
}

func (b *Redis) HGetAllMap(ctx context.Context, key string) (map[string]string, error) {
	data, err := b.client.HGetAll(ctx, b.PrefixKey(key)).Result()
	if len(data) == 0 {
		return nil, ErrNoData
	}
	return data, err
}

// HMSet support gosql.Model and map[string]interface{}
func (b *Redis) HMSet(ctx context.Context, key string, m any) error {
	return b.hMSet(ctx, b.PrefixKey(key), m, CACHE_HOUR_TTL, ScanTagName)
}

// HMSetWithExpire support set expire time
func (b *Redis) HMSetWithExpire(ctx context.Context, key string, m any, expire time.Duration) error {
	return b.hMSet(ctx, b.PrefixKey(key), m, expire, ScanTagName)
}

func (b *Redis) hMSet(ctx context.Context, key string, m any, expire time.Duration, tagName string) error {
	var mm map[string]any
	var ok bool
	ref := reflect.TypeOf(m)
	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}
	switch ref.Kind() {
	case reflect.Struct:
		mm = structToMap(m, tagName)
	case reflect.Map:
		mm, ok = m.(map[string]any)
		if !ok {
			ms, ok := m.(map[string]string)
			if !ok {
				return errors.New("value must is map[string]interafce{}")
			}
			mm = make(map[string]any)
			for key, value := range ms {
				mm[key] = value
			}
		}
	default:
		return fmt.Errorf("cannot convert from %s", ref)
	}

	_, err := b.client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(ctx, key, mm)
		pipe.Expire(ctx, key, expire)
		return nil
	})
	return err
}

func (b *Redis) HGet(ctx context.Context, key string, field string) (string, error) {
	return b.client.HGet(ctx, b.PrefixKey(key), field).Result()
}

func (b *Redis) HSet(ctx context.Context, key string, field string, value any) error {
	return b.client.HSet(ctx, b.PrefixKey(key), field, value).Err()
}

func (b *Redis) HDel(ctx context.Context, key string, field string) error {
	return b.client.HDel(ctx, b.PrefixKey(key), field).Err()
}

func (b *Redis) HExists(ctx context.Context, key string, field string) (bool, error) {
	return b.client.HExists(ctx, b.PrefixKey(key), field).Result()
}

func (b *Redis) HIncrBy(ctx context.Context, key string, field string, incr int64) (int64, error) {
	return b.client.HIncrBy(ctx, b.PrefixKey(key), field, incr).Result()
}

func (b *Redis) Exists(ctx context.Context, key string) bool {
	n, err := b.client.Exists(ctx, b.PrefixKey(key)).Result()

	if err != nil || n != 1 {
		return false
	}

	return true
}

func (b *Redis) LPop(ctx context.Context, key string) (string, error) {
	return b.client.LPop(ctx, b.PrefixKey(key)).Result()
}

func (b *Redis) LLen(ctx context.Context, key string) (int64, error) {
	return b.client.LLen(ctx, b.PrefixKey(key)).Result()
}

func (b *Redis) RPop(ctx context.Context, key string) (string, error) {
	return b.client.RPop(ctx, b.PrefixKey(key)).Result()
}

func (b *Redis) LPush(ctx context.Context, key string, values ...any) (int64, error) {
	return b.client.LPush(ctx, b.PrefixKey(key), values...).Result()
}

func (b *Redis) RPush(ctx context.Context, key string, values ...any) (int64, error) {
	return b.client.RPush(ctx, b.PrefixKey(key), values...).Result()
}

func (b *Redis) SAdd(ctx context.Context, key string, members ...any) (int64, error) {
	return b.client.SAdd(ctx, b.PrefixKey(key), members...).Result()
}

func (b *Redis) TTL(ctx context.Context, key string) time.Duration {
	return b.client.TTL(ctx, b.PrefixKey(key)).Val()
}

func (b *Redis) GeoAdd(ctx context.Context, key string, name string, lng float64, lat float64) error {
	return b.client.GeoAdd(ctx, b.PrefixKey(key), &redis.GeoLocation{
		Name:      name,
		Longitude: lng,
		Latitude:  lat,
	}).Err()
}

func (b *Redis) GeoRadius(ctx context.Context, key string, lng float64, lat float64, radius float64, unit string) ([]redis.GeoLocation, error) {

	return b.client.GeoRadius(ctx, b.PrefixKey(key), lng, lat, &redis.GeoRadiusQuery{
		Radius:   radius,
		Unit:     unit,
		Sort:     "ASC",
		WithDist: true,
	}).Result()

}

func (b *Redis) Publish(ctx context.Context, channel string, message any) error {
	return b.client.Publish(ctx, b.PrefixKey(channel), message).Err()
}

func (b *Redis) Subscribe(ctx context.Context, channels ...string) *PubSub {
	chans := make([]string, len(channels))

	for i, channel := range channels {
		chans[i] = b.PrefixKey(channel)
	}

	return b.client.Subscribe(ctx, chans...)
}

type Pipeliner = redis.Pipeliner

func (b *Redis) Multi(ctx context.Context, fn func(pipe Pipeliner) error) error {
	_, err := b.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		return fn(pipe)
	})
	return err
}

var ErrNoData = errors.New("data is empty")
var ErrClosed = redis.ErrClosed

func IsNil(err error) bool {
	return errors.Is(err, redis.Nil)
}

const Nil = redis.Nil
