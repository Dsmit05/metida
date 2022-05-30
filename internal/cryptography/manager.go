package cryptography

import "time"

type ManagerToken interface {
	AccessToken
	RefreshToken
}

type ManagerToken1 interface {
	CreateToken(email string, role string, ttl time.Duration) (string, error)
	RefreshToken
}

func NewManagerToken(secret string) ManagerToken {
	accessToken := NewTokenJWT(secret)
	refreshToken := NewRefreshToken()

	return struct {
		AccessToken
		RefreshToken
	}{accessToken, refreshToken}
}
