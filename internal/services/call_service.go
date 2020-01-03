package services

import "webrtc-server/driver"

import "webrtc-server/internal/repositories"

// CallService struct
type CallService struct {
	db *driver.Database
}

// RegisterSocketID func
func (c *CallService) RegisterSocketID(token string) bool {
	return true
}

// NewCallService func
func NewCallService(db *driver.Database) repositories.CallRepository {
	return &CallService{
		db: db,
	}
}
