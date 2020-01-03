package jwtauth

import (
	"fmt"
	"webrtc-server/internal/models"

	"github.com/dgrijalva/jwt-go"
)

// ParseTokenToUser ...
func ParseTokenToUser(tokenStr string) (models.User, error) {
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

	// db := models.OpenDB()
	// var baseUser BaseUser
	// mapstructure.Decode(token.Claims, &baseUser)
	// id, _ := strconv.ParseUint(fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["id"]), 10, 32)
	// baseUser.ID = uint(id)
	// baseUser.NationalNumber = token.Claims.(jwt.MapClaims)["national_number"].(string)
	// baseUser.CountryPrefix = token.Claims.(jwt.MapClaims)["country_prefix"].(string)
	// baseUser.FirstName = token.Claims.(jwt.MapClaims)["first_name"].(string)
	// baseUser.LastName = token.Claims.(jwt.MapClaims)["last_name"].(string)
	// baseUser.Birthday = token.Claims.(jwt.MapClaims)["birthday"].(string)
	// baseUser.CountryCode = token.Claims.(jwt.MapClaims)["country_code"].(string)
	// baseUser.City = token.Claims.(jwt.MapClaims)["city"].(string)
	// baseUser.Gender = token.Claims.(jwt.MapClaims)["gender"].(string)

	// db.Where("id = ? AND national_number = ? and country_prefix = ?", baseUser.ID, baseUser.NationalNumber, baseUser.CountryPrefix).First(&user)

	// if user.ID == 0 {
	// 	return user, fmt.Errorf("Invalid authorization token - Does not match UserID")
	// }

	return user, nil
}
