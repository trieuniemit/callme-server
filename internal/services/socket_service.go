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
func (c *SocketService) RegisterSocketID(user *models.User) ([]string, error) {
	err := c.db.Conn.Model(user).UpdateColumn("socket_id", gorm.Expr("?", user.SocketID)).Error
	socketIDs := []string{}

	if err != nil {
		return socketIDs, err
	}

	subQuery := c.db.Conn.Model(models.Contact{}).Select("user2").Where("user1=?", user.ID).SubQuery()
	rows, err := c.db.Conn.Raw("SELECT socket_id FROM users WHERE (id IN (?) AND socket_id IS NOT NULL)", subQuery).Rows()

	if err != nil {
		return socketIDs, err
	}

	for rows.Next() {
		ID := ""
		rows.Scan(&ID)
		socketIDs = append(socketIDs, ID)
	}

	return socketIDs, err
}

// SetCallingStatus func
func (c *SocketService) SetCallingStatus(status bool, IDs []uint) {
	c.db.Conn.Exec("UPDATE users SET calling=? WHERE id IN (?)", status, IDs)
}

// RemoveSocketIDs func
func (c *SocketService) RemoveSocketIDs(IDs []uint) {
	c.db.Conn.Exec("UPDATE users SET socket_id=NULL WHERE id IN (?)", IDs)
}

// NewSocketService func
func NewSocketService(db *driver.Database) repositories.SocketRepository {
	return &SocketService{
		db: db,
	}
}
