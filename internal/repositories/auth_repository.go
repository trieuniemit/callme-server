package repositories

import (
	"webrtc-server/internal/models"
)

// AuthRepository ...
type AuthRepository interface {
	Register(user *models.User) *models.User
	Login(user *models.User) (*models.User, error)
	Logout() error
}
