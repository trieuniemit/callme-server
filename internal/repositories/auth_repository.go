package repositories

import (
	"webrtc-server/internal/models"
)

// AuthRepository ...
type AuthRepository interface {
	Register(user *models.User) *models.User
	Login(username string) (models.User, error)
	Logout() error
}
