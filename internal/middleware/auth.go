package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"webrtc-server/internal/handler/response"

	"github.com/dgrijalva/jwt-go"
)

// Auth validator auth req
func (m *Middleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		// Kiểm tra xem có tồn tại token không
		if authorizationHeader == "" {
			resData := response.Message(false, "An authorization header is required")
			resData["key"] = "invalid_token"
			response.RespondUnauthorized(w, resData)
			return
		}
		// Kiểm tra xem token có đúng định dạng Bearer token không
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) != 2 {
			resData := response.Message(false, "Invalid authorization token")
			resData["key"] = "invalid_token"
			response.RespondUnauthorized(w, resData)
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			resData := response.Message(false, err.Error())
			resData["key"] = "invalid_token"
			response.RespondInternalServer(w, resData)
			return
		}
		if !token.Valid {
			resData := response.Message(false, "Invalid authorization token")
			resData["key"] = "invalid_token"
			response.RespondUnauthorized(w, resData)
			return
		}

		//id, _ := strconresponse.ParseUint(fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["id"]), 10, 32)

		log.Println(token.Claims.(jwt.MapClaims))

		// db.Where("id = ? AND national_number = ? and country_prefix = ?", baseUser.ID, baseUser.NationalNumber, baseUser.CountryPrefix).First(&user)
		// if user.ID == 0 {
		// 	resData := response.Message(false, "Invalid authorization token - Does not match UserID")
		// 	resData["key"] = "invalid_token"
		// 	response.RespondUnauthorized(w, resData)
		// 	return
		// }

		// ctx := context.WithValue(r.Context(), "user", user)
		// r = r.WithContext(ctx)
		next(w, r)
	})
}
