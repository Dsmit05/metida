package cryptography

import (
	"fmt"

	"github.com/google/uuid"
)

type RefreshToken interface {
	CreateRefreshToken() (string, error)
}

type TokenRefresh struct{}

func NewRefreshToken() *TokenRefresh {
	return &TokenRefresh{}
}

func (o *TokenRefresh) CreateRefreshToken() (string, error) {
	rToken := uuid.New()
	if rToken.String() == "" {
		return "", fmt.Errorf("Token is empty")
	}

	return rToken.String(), nil
}
