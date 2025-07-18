package repository

import (
	"context"
	"inkgo/database"
	"time"
)

type TokenRepository interface {
	// SetToken 设置一次性token
	SetToken(token string, value string, expiration time.Duration) error
	// GetToken 获取一次性token
	GetToken(token string) (string, error)
}

type tokenRepository struct {
	rdb *database.RedisDB
}

func NewTokenRepository(rdb *database.RedisDB) TokenRepository {
	return &tokenRepository{
		rdb: rdb,
	}
}
func (t *tokenRepository) SetToken(token string, value string, expiration time.Duration) error {
	return t.rdb.Set(context.Background(), token, value, expiration)
}

func (t *tokenRepository) GetToken(token string) (string, error) {
	return t.rdb.Get(context.Background(), token)
}
