package repositories

import "webrtc-server/internal/models"

// ContactRepository interface
type ContactRepository interface {
	GetList(user *models.User) ([]models.User, error)
	AddContact(user *models.User, userID uint) error
}
