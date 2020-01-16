package repositories

import "webrtc-server/internal/models"

// SocketRepository interface
type SocketRepository interface {
	RegisterSocketID(user *models.User) ([]string, error)
	SetCallingStatus(status bool, IDs []uint)
	RemoveSocketIDs(IDs []uint)
}
