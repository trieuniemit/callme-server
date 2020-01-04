package jwtauth

import (
	"fmt"
	"webrtc-server/driver"
	"webrtc-server/internal/models"

	"github.com/dgrijalva/jwt-go"
)

// ParseTokenToUser ...
func ParseTokenToUser(tokenStr string, db *driver.Database) (models.User, error) {
	var user models.User
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return user, fmt.Errorf("There was an error")
	}
	if !token.Valid {
		return user, fmt.Errorf("Invalid authorization token")
	}

	userMap := (token.Claims.(jwt.MapClaims))

	db.Conn.Where("password = ? AND email = ?", userMap["password"], userMap["email"]).First(&user)
	if user.ID == 0 {
		return user, fmt.Errorf("User not found")
	}

	return user, nil
}
