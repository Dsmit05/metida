package cryptography

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type AccessToken interface {
	CreateToken(email string, role string, ttl time.Duration) (string, error)
	ParseToken(inputToken string) (email string, role string, err error)
}

type TokenJWT struct {
	secretKey []byte
}

func NewTokenJWT(secret string) *TokenJWT {
	return &TokenJWT{secretKey: []byte(secret)}
}

// UserClaims include custom claims on jwt.
type UserClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

// CreateToken create new token with parameters.
func (o *TokenJWT) CreateToken(email, role string, ttl time.Duration) (string, error) {
	claims := UserClaims{email, role,
		jwt.StandardClaims{ExpiresAt: time.Now().Add(ttl).Unix()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(o.secretKey)
}

// ParseToken parsing input token, and return email and role from token.
func (o *TokenJWT) ParseToken(inputToken string) (email, role string, err error) {
	token, err := jwt.Parse(inputToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return o.secretKey, nil
	})
	if err != nil {
		return "", "", err
	}

	if !token.Valid {
		return "", "", fmt.Errorf("not valid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("error get user claims from token")
	}

	return claims["email"].(string), claims["role"].(string), nil
}
