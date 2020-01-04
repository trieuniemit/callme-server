package middleware

import (
	"context"
	"net/http"
	"strings"
	"webrtc-server/driver"
	"webrtc-server/internal/handler/response"
	"webrtc-server/pkg/jwtauth"
)

// Auth validator auth req
func Authenticate(next http.HandlerFunc, db *driver.Database) http.HandlerFunc {
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

		user, err := jwtauth.ParseTokenToUser(bearerToken[1], db)
		if err != nil {
			resData := response.Message(false, err.Error())
			resData["key"] = "invalid_token"
			response.RespondInternalServer(w, resData)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)
		next(w, r)
	})
}
