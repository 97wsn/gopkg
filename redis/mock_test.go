package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTestRedis(t *testing.T) {
	rds := NewTestRedis()
	ctx := context.Background()
	assert.False(t, rds.Exists(ctx, "key"))
	assert.NoError(t, rds.Set(ctx, "key", "val"))
	val, err := rds.Get(ctx, "key")
	assert.NoError(t, err)
	assert.Equal(t, "val", val)
	assert.NoError(t, rds.Unlink(ctx, "key"))
	assert.False(t, rds.Exists(ctx, "key"))
}
