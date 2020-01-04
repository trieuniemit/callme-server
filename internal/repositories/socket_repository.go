package repositories

// SocketRepository interface
type SocketRepository interface {
	RegisterSocketID(token string) bool
	SetCallingStatus(status bool, IDs []uint)
}
