package jwtauth

import (
	"time"
	"webrtc-server/internal/models"

	"github.com/dgrijalva/jwt-go"
)

// TokenClaim struct
type TokenClaim struct {
	*jwt.StandardClaims
	models.User
}

// CreateToken ...
func CreateToken(user models.User) (string, int64, error) {
	expiresAt := time.Now().Add(time.Hour * 24 * 30).Unix()

	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = &TokenClaim{
		&jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		user,
	}

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		return "", 0, error
	}
	return tokenString, expiresAt, nil
}
