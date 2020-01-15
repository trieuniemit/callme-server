package services

import (
	"webrtc-server/driver"
	"webrtc-server/internal/models"
	"webrtc-server/internal/repositories"

	"github.com/jinzhu/gorm"
)

// SocketService struct
type SocketService struct {
	db *driver.Database
}

// RegisterSocketID func
func (c *SocketService) RegisterSocketID(user *models.User) error {
	return c.db.Conn.Model(user).UpdateColumn("socket_id", gorm.Expr("?", user.SocketID)).Error
}

// SetCallingStatus func
func (c *SocketService) SetCallingStatus(status bool, IDs []uint) {
	c.db.Conn.Exec("UPDATE users SET calling=? WHERE id IN (?)", status, IDs)
}

// NewSocketService func
func NewSocketService(db *driver.Database) repositories.SocketRepository {
	return &SocketService{
		db: db,
	}
}
