package services

import (
	"webrtc-server/driver"
	"webrtc-server/internal/models"
	"webrtc-server/internal/repositories"
)

// NewUserService retunrs implement of post repository interface
func NewUserService(db *driver.Database) repositories.UserRepository {
	return &userImpl{
		database: db,
	}
}

type userImpl struct {
	database *driver.Database
}

func (instance *userImpl) List(num int64) ([]*models.User, error) {
	var users []*models.User
	instance.database.Conn.Limit(num).Find(&users)
	return users, nil
}

func (instance *userImpl) GetByID(id int64) (*models.User, error) {
	var user *models.User
	if err := instance.database.Conn.First(&user, id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (instance *userImpl) Create(t *models.User) (int64, error) {
	instance.database.Conn.Create(t)
	return int64(t.ID), nil
}

func (instance *userImpl) Update(t *models.User) (*models.User, error) {
	instance.database.Conn.Update(t)
	return t, nil
}

func (instance *userImpl) Delete(id int64) (bool, error) {
	var user *models.User
	if err := instance.database.Conn.First(&user, id).Error; err != nil {
		return false, err
	}
	instance.database.Conn.Delete(&user)
	return true, nil
}
