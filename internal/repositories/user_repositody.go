package repositories

import (
	"webrtc-server/internal/models"
)

// UserRepository ...
type UserRepository interface {
	List(num int64) ([]*models.User, error)
	GetByID(id int64) (*models.User, error)
	Create(user *models.User) (int64, error)
	Update(user *models.User) (*models.User, error)
	Delete(id int64) (bool, error)
}
