package services

import (
	"log"
	"webrtc-server/driver"
	"webrtc-server/internal/models"
	"webrtc-server/internal/repositories"
)

type userImpl struct {
	database *driver.Database
}

// NewUserService retunrs implement of post repository interface
func NewUserService(db *driver.Database) repositories.UserRepository {
	return &userImpl{
		database: db,
	}
}

func (instance *userImpl) List(num int64) ([]*models.User, error) {
	var users []*models.User
	instance.database.Conn.Limit(num).Find(&users)
	return users, nil
}

func (instance *userImpl) GetByID(id int64) (*models.User, error) {
	user := models.User{}

	if err := instance.database.Conn.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (instance *userImpl) Create(u *models.User) (int64, error) {
	instance.database.Conn.Create(u)
	log.Println(u)
	return int64(u.ID), nil
}

func (instance *userImpl) Update(u *models.User) (*models.User, error) {
	instance.database.Conn.Debug().Model(u).Update(u)
	return u, nil
}

func (instance *userImpl) Delete(id int64) (bool, error) {
	var user *models.User
	if err := instance.database.Conn.First(&user, id).Error; err != nil {
		return false, err
	}
	instance.database.Conn.Delete(&user)
	return true, nil
}
