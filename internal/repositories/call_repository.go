package repositories

// CallRepository interface
type CallRepository interface {
	RegisterSocketID(token string) bool
}
