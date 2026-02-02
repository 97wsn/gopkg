package redis

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getConn(t *testing.T) *redis.Client {
	conf := &Config{
		Server:       redisServer,
		MaxRetries:   3,
		DialTimeout:  2,
		ReadTimeout:  3,
		WriteTimeout: 4,
	}

	client, err := Connect(conf)
	require.NoError(t, err)
	client.FlushDB(context.Background())
	return client
}

func TestRedis(t *testing.T) {
	t.Skip("skip this test for now")
	client := NewRedis(getConn(t))
	ctx := context.Background()
	assert.False(t, client.Exists(ctx, "key"))
	assert.NoError(t, client.Set(ctx, "key", "val"))
	val, err := client.Get(ctx, "key")
	assert.NoError(t, err)
	assert.Equal(t, "val", val)
	assert.NoError(t, client.Unlink(ctx, "key"))
	assert.False(t, client.Exists(ctx, "key"))

}
