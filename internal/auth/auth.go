package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator interface {
	GenerateToken(userID int64, ttl time.Duration) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
