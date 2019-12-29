package repositories

import (
	"webrtc-server/internal/models"
)

// AuthRepository ...
type AuthRepository interface {
	Register(user *models.User) *models.User
	Login(email string) *models.User
	Logout() error
}
