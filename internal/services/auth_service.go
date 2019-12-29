package services

import (
	"webrtc-server/driver"
	"webrtc-server/internal/models"
	"webrtc-server/internal/repositories"
)

type authImpl struct {
	db *driver.Database
}

func (auth *authImpl) Login(user *models.User) (*models.User, error) {
	return nil, nil
}

func (auth *authImpl) Register(user *models.User) *models.User {
	if err := auth.db.Conn.Create(user).Error; err == nil {
		return user
	}
	return nil
}

func (auth *authImpl) Logout() error {
	return nil
}

// NewAuthService ...
func NewAuthService(db *driver.Database) repositories.AuthRepository {
	return &authImpl{
		db: db,
	}
}
