package loop

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoop(t *testing.T) {
	assert.NoError(t, Loop[int](
		context.Background(),
		1,
		50,
		BackoffConfig{
			MaxRetries: 5,
			BaseDelay:  30 * time.Millisecond,
			MaxDelay:   2 * time.Second,
		},
		30*time.Millisecond,
		func(ctx context.Context, page int, pageSize int) (int, bool, error) {
			fmt.Printf("Processing page %d (size %d)\n", page, pageSize)
			// 判断是否结束
			isEnd := page >= 5
			return page + 1, isEnd, nil
		}))
}
