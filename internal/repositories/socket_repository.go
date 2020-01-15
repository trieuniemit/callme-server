package repositories

import "webrtc-server/internal/models"

// SocketRepository interface
type SocketRepository interface {
	RegisterSocketID(user *models.User) error
	SetCallingStatus(status bool, IDs []uint)
}
