package redis

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
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
