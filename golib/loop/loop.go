package loop

import (
	"context"
	"fmt"
	"math"
	"time"
)

// Handler 处理分页数据的回调函数类型
type Handler[T int | string] func(ctx context.Context, page T, pageSize int) (T, bool, error)

// BackoffConfig 指数退避配置
type BackoffConfig struct {
	MaxRetries int           // 最大重试次数
	BaseDelay  time.Duration // 基础延迟时间
	MaxDelay   time.Duration // 最大延迟时间
}

// 指数退避默认配置
var defaultBackoff = BackoffConfig{
	MaxRetries: 3,
	BaseDelay:  100 * time.Millisecond,
	MaxDelay:   5 * time.Second,
}

// Loop 通用分页处理循环，支持指数退避重试和上下文取消
func Loop[T int | string](ctx context.Context, page T, pageSize int,
	backoff BackoffConfig, rateLimit time.Duration,
	fn Handler[T]) error {
	// 使用默认配置（如果未提供）
	if backoff.MaxRetries <= 0 {
		backoff = defaultBackoff
	}

	// 检查上下文是否已取消
	if ctx.Err() != nil {
		return ctx.Err()
	}

	for {
		// 执行回调函数
		nextPage, isEnd, err := fn(ctx, page, pageSize)
		if err != nil {
			// 带指数退避的重试机制
			for retry := 0; retry < backoff.MaxRetries; retry++ {
				// 检查上下文是否已取消
				if ctx.Err() != nil {
					return ctx.Err()
				}

				// 计算退避延迟
				delay := time.Duration(math.Pow(2, float64(retry))) * backoff.BaseDelay
				if delay > backoff.MaxDelay {
					delay = backoff.MaxDelay
				}

				// 使用带超时的睡眠，支持提前取消
				timer := time.NewTimer(delay)
				select {
				case <-ctx.Done():
					timer.Stop()
					return ctx.Err()
				case <-timer.C:
				}

				// 重试回调
				if nextPage, isEnd, err = fn(ctx, page, pageSize); err == nil {
					break
				}
			}

			// 所有重试都失败
			if err != nil {
				return fmt.Errorf("callback failed after %d retries: %w", backoff.MaxRetries, err)
			}
		}

		// 检查是否结束
		if isEnd {
			return nil
		}

		// 更新页码
		page = nextPage

		// 控制请求频率（支持上下文取消）
		if rateLimit > 0 {
			timer := time.NewTimer(rateLimit)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}
	}
}
