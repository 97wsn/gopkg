package jwtutil

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type KeyFunc func(ctx context.Context) jwt.Keyfunc
type DataFunc func() any

type Claims struct {
	jwt.RegisteredClaims
	Data any `json:"data,omitempty"`
}

func GenerateToken(data any, key []byte, expire time.Duration) (string, error) {
	claims := Claims{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
}

func ParseToken(token string, keyFunc jwt.Keyfunc, dataFunc DataFunc) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{
		Data: dataFunc(),
	}, keyFunc)
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, errors.New("invalid token")
}

func ParseWith[T any](token string, keyFunc jwt.Keyfunc) (*T, error) {
	var m T
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{
		Data: &m,
	}, keyFunc)
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims.Data.(*T), nil
		}
	}
	return nil, errors.New("invalid token")
}
