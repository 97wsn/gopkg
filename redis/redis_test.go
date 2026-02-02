package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var redisServer = "localhost:6379"

func TestConnect(t *testing.T) {
	t.Skip("skip this test for now")
	client, err := Connect(&Config{Server: redisServer})
	assert.NoError(t, err)
	assert.NotEqual(t, client, nil, "Redis connect error")
}

func TestConnectError(t *testing.T) {
	conf := &Config{
		Server: "111.123.111.111",
	}

	client, err := Connect(conf)
	assert.Error(t, err)
	assert.NotEqual(t, client, nil, "Redis connect error")
}
