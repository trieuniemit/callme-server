package middleware

import (
	"webrtc-server/driver"
)

// Middleware struct
type Middleware struct {
	db *driver.Database
}

// NewMiddleware ...
func NewMiddleware(db *driver.Database) *Middleware {
	return &Middleware{
		db: db,
	}
}
