package service

import (
	"inkgo/repository"
	"time"
)

type TokenService interface {
	// SetToken 设置一次性token
	SetToken(token string, value string, expiration time.Duration) error
	// GetToken 获取一次性token
	GetToken(token string) (string, error)
}

type tokenService struct {
	repo repository.TokenRepository
}

func NewTokenService(repo repository.TokenRepository) TokenService {
	return &tokenService{
		repo: repo,
	}
}

func (t *tokenService) SetToken(token string, value string, expiration time.Duration) error {
	return t.repo.SetToken(token, value, expiration)
}

func (t *tokenService) GetToken(token string) (string, error) {
	return t.repo.GetToken(token)
}
