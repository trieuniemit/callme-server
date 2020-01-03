package jwtauth

import (
	"net/http"
	"strings"
)

// GetToken func
func GetToken(r *http.Request) string {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return ""
	}
	bearerToken := strings.Split(authorizationHeader, " ")
	if len(bearerToken) != 2 {
		return ""
	}
	return bearerToken[1]
}
