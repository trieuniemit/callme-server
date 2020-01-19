package services

import (
	"webrtc-server/driver"
	"webrtc-server/internal/models"
	"webrtc-server/internal/repositories"
)

type contactService struct {
	db *driver.Database
}

func (c *contactService) GetList(user *models.User) ([]models.User, error) {
	users := []models.User{}

	listQuery := c.db.Conn.Model(&models.Contact{}).Select("user2").Where("user1 = ?", user.ID).SubQuery()
	err := c.db.Conn.Model(user).Where("id IN(?)", listQuery).Find(&users).Error
	return users, err
}

func (c *contactService) AddContact(user *models.User, userID uint) error {
	contact := models.Contact{}
	contact.User1 = user.ID
	contact.User2 = userID
	return c.db.Conn.Create(contact).Error
}

// NewContactService ...
func NewContactService(db *driver.Database) repositories.ContactRepository {
	return &contactService{
		db: db,
	}
}
