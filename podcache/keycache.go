package podcache

import (
	"context"
	"database/sql"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

var (
	maxSize              = 1000
	expireAfterInSeconds = 60 * 2
)

type KeyCache[K comparable, V any] interface {
	Get(ctx context.Context, key K) (V, error)
	BatchGet(ctx context.Context, keys []K) (map[K]V, error)
}

type Option[K comparable, V any] func(b *keyCache[K, V])

// WithMaxSize 设置缓存最大数量
func WithMaxSize[K comparable, V any](maxSize int) Option[K, V] {
	return func(b *keyCache[K, V]) {
		if maxSize > 0 {
			b.maxSize = maxSize
		}
	}
}

// WithExpires 设置缓存过期时间
func WithExpires[K comparable, V any](expires int) Option[K, V] {
	return func(b *keyCache[K, V]) {
		if expires > 0 {
			b.expireAfterSeconds = expires
		}
	}
}

type keyData[V any] struct {
	createdAt time.Time
	data      V
}

// NewKeyCache 创建一个新的POD缓存
func NewKeyCache[K comparable, V any](mutex *sync.Mutex, load func(ctx context.Context, keys []K) (map[K]V, error), options ...Option[K, V]) KeyCache[K, V] {
	p := &keyCache[K, V]{
		load:               load,
		mutex:              mutex,
		maxSize:            maxSize,
		expireAfterSeconds: expireAfterInSeconds,
	}

	if len(options) > 0 {
		for _, option := range options {
			option(p)
		}
	}

	p.lru, _ = lru.New[K, keyData[V]](p.maxSize)

	return p
}

type keyCache[K comparable, V any] struct {
	// 回源机制
	load func(ctx context.Context, keys []K) (map[K]V, error)
	// 缓存有效期，超过该时间则重新回源
	expireAfterSeconds int
	// 缓存数据
	lru *lru.Cache[K, keyData[V]]
	// 缓存的最大数量
	maxSize int
	mutex   *sync.Mutex
}

func (p *keyCache[K, V]) Get(ctx context.Context, key K) (V, error) {
	v, ok := p.getFromLru(key)
	if ok {
		return v, nil
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 防止并发重复回源
	v, ok = p.getFromLru(key)
	if ok {
		return v, nil
	}

	// 回源数据失败，直接返回错误，由调用方处理
	result, err := p.load(ctx, []K{key})
	if err != nil {
		var none V
		return none, err
	}
	var none V

	if len(result) == 0 {
		return none, sql.ErrNoRows
	}

	data, ok := result[key]
	if !ok {
		return none, sql.ErrNoRows
	}

	// 设置lru缓存
	p.lru.Add(key, keyData[V]{
		createdAt: time.Now(),
		data:      data,
	})

	return data, nil
}

func (p *keyCache[K, V]) BatchGet(ctx context.Context, keys []K) (map[K]V, error) {

	exists, nonExists := make(map[K]V), make(map[K]V)

	for _, key := range keys {
		v, ok := p.getFromLru(key)
		if ok {
			exists[key] = v
		} else {
			nonExists[key] = v
		}
	}

	// 全部查到了，无需回源
	if len(nonExists) == 0 {
		return exists, nil
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	againNonExists := make([]K, 0)

	for key := range nonExists {
		v, ok := p.getFromLru(key)
		if ok {
			exists[key] = v
		} else {
			againNonExists = append(againNonExists, key)
		}
	}

	// 回源数据失败，直接返回错误，由调用方处理
	result, err := p.load(ctx, againNonExists)
	if err != nil {
		var none map[K]V
		return none, err
	}

	// 设置lru缓存
	for key, value := range result {
		p.lru.Add(key, keyData[V]{
			createdAt: time.Now(),
			data:      value,
		})
		// 本次回源后要返回的结果
		exists[key] = value
	}

	return exists, nil
}

func (p *keyCache[K, V]) getFromLru(key K) (V, bool) {
	v, ok := p.lru.Get(key)
	// 当前缓存存在数据
	if ok && v.createdAt.Add(time.Duration(p.expireAfterSeconds)*time.Second).After(time.Now()) {

		return v.data, true
	}
	var none V
	return none, false
}
