package jwtutil

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestParseToken(t *testing.T) {
	jwtKey := "e4da3b7fbbce2345d7772b0674a318d5"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njk2Njk3MTIsImRhdGEiOnsidXNlcl9pZCI6MSwiQXV0aFR5cGUiOiJBRE1JTiJ9fQ.iuELZmqytKj8al46CkLx-c1t4CCyW1HgX_9kQHeZGio"

	type Data struct {
		UserId int32 `json:"user_id"`
	}

	cls, err := ParseToken(token, func(token *jwt.Token) (any, error) {
		return []byte(jwtKey), nil
	}, func() any {
		return &Data{}
	})

	assert.True(t, errors.Is(err, jwt.ErrTokenExpired))
	assert.Nil(t, cls)
}

func TestParseWith(t *testing.T) {
	jwtKey := "e4da3b7fbbce2345d7772b0674a318d5"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjUzNjk2NjYwMjksImRhdGEiOnsidXNlcl9pZCI6MSwiQXV0aFR5cGUiOiJBRE1JTiJ9fQ.IL3NmMvKVZRNBozlspHyExCEi7FxOO-UHEYbmhGCnk0"

	type Data struct {
		UserId int32 `json:"user_id"`
	}

	data, err := ParseWith[Data](token, func(token *jwt.Token) (any, error) {
		return []byte(jwtKey), nil
	})

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, int32(1), data.UserId)
}

func TestJWT(t *testing.T) {
	if testing.Short() {
		return
	}

	type Data struct {
		UserId   int32 `json:"user_id"`
		AuthType string
	}

	t.Run("jwt", func(t *testing.T) {
		token, err := GenerateToken(Data{
			UserId:   1,
			AuthType: "ADMIN",
		}, []byte("e4da3b7fbbce2345d7772b0674a318d5"), 10000*time.Hour)
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	})
}
