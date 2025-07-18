package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"inkgo/database"
	"time"
)

type AuthRepository interface {
	SaveEmailCode(ctx context.Context, email, code string, expiration time.Duration) error
	FindEmailCode(ctx context.Context, email string) (string, error)
	SaveResetToken(ctx context.Context, token, email string, expiration time.Duration) error
	FindEmailByResetToken(ctx context.Context, token string) (string, error)
}

type authRepository struct {
	rdb *database.RedisDB
}

func NewAuthRepository(rdb *database.RedisDB) AuthRepository {
	return &authRepository{
		rdb: rdb,
	}
}

func (e *authRepository) SaveEmailCode(ctx context.Context, email, code string, expiration time.Duration) error {
	key := generateVerifyCodeKey(email)
	err := e.rdb.Set(ctx, key, code, expiration)
	if err != nil {
		return fmt.Errorf("failed to set verification code: %w", err)
	}
	return nil
}

func (e *authRepository) FindEmailCode(ctx context.Context, email string) (string, error) {
	key := generateVerifyCodeKey(email)
	result, err := e.rdb.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("verification code not found for email: %s", email)
		}
		return "", fmt.Errorf("failed to get verification code: %w", err)
	}
	return result, nil
}

func generateVerifyCodeKey(email string) string {
	return fmt.Sprintf("mail_%s", email)
}

func (e *authRepository) SaveResetToken(ctx context.Context, token, email string, expiration time.Duration) error {
	key := fmt.Sprintf("reset_token:%s", token)
	if err := e.rdb.Set(ctx, key, email, expiration); err != nil {
		return err
	}
	return nil
}

func (e *authRepository) FindEmailByResetToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("reset_token:%s", token)
	email, err := e.rdb.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("reset token not found: %s", token)
		}
		return "", fmt.Errorf("failed to get email by reset token: %w", err)
	}
	return email, nil
}
