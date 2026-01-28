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
	jwtKey := "abcde123456789"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDc1Njg1NDMsIkRhdGEiOnsibWVyY2hhbnRfaWQiOjF9fQ.4z9DpP6PFPgIVwYYU8hGaIuudjPdGt4ok7EfOkKVqs4"

	type Data struct {
		MerchantId int32 `json:"merchant_id"`
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
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjUzMTM0MTc1NTMsImRhdGEiOnsibWVyY2hhbnRfaWQiOjEsIkF1dGhUeXBlIjoiQURNSU4ifX0.rH0XbXnNWbYLw2MNPviOvBSbf_81lIroY5l9khfA8tM"

	type Data struct {
		MerchantId int32 `json:"merchant_id"`
	}

	data, err := ParseWith[Data](token, func(token *jwt.Token) (any, error) {
		return []byte(jwtKey), nil
	})

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, int32(1), data.MerchantId)
}

func TestJWT(t *testing.T) {
	if testing.Short() {
		return
	}

	type Data struct {
		MerchantId int32 `json:"merchant_id"`
		AuthType   string
	}

	t.Run("jwt", func(t *testing.T) {
		token, err := GenerateToken(Data{
			MerchantId: 1,
			AuthType:   "ADMIN",
		}, []byte("e4da3b7fbbce2345d7772b0674a318d5"), 999999*time.Hour)
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	})
}
