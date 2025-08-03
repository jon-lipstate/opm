package helpers

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"opm/middleware"
)

// GenerateJWT creates a new JWT token for a user
func GenerateJWT(userID int, secret string) (string, error) {
	claims := &middleware.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}