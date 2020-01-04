package services

import (
	"webrtc-server/driver"
	"webrtc-server/internal/repositories"
)

// SocketService struct
type SocketService struct {
	db *driver.Database
}

// RegisterSocketID func
func (c *SocketService) RegisterSocketID(token string) bool {
	return true
}

// SetCallingStatus func
func (c *SocketService) SetCallingStatus(status bool, IDs []uint) {
	c.db.Conn.Debug().Exec("UPDATE users SET calling=? WHERE id IN (?)", status, IDs)
}

// NewSocketService func
func NewSocketService(db *driver.Database) repositories.SocketRepository {
	return &SocketService{
		db: db,
	}
}
