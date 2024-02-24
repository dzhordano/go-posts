package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenIsExpired = errors.New("token has expired")
	ErrTokenIsInvalid = errors.New("token is invalid")
)

type TokenManager interface {
	CreateJWT(userid string, ttl time.Duration) (string, error)
	ValidateToken(refreshToken string) (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if len(signingKey) < 16 {
		return nil, errors.New("signing key length must be at least 16 characters")
	}

	return &Manager{
		signingKey: signingKey,
	}, nil
}

func (m *Manager) CreateJWT(userId string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) ValidateToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, ErrTokenIsInvalid
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrTokenIsInvalid
	}

	return claims["sub"].(string), nil
}
