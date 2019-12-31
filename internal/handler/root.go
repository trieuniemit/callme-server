package handler

import (
	"webrtc-server/driver"
	"webrtc-server/internal/handler/http"
	"webrtc-server/internal/handler/socket"

	"github.com/gorilla/mux"
)

// RegisterHTTP ...
func RegisterHTTP(db *driver.Database, r *mux.Router) {
	apiPrefix := r.PathPrefix("/api").Subrouter()
	registerAPI(db, apiPrefix)

	// register websocket
	socket.RegisterSocketRoute(socket.NewSocketHandler(db), r)
}

func registerAPI(db *driver.Database, r *mux.Router) {
	// register api
	http.RegisterAuthRoutes(http.NewAuthHandler(db), r)
	http.RegisterUserRoutes(http.NewUserHandler(db), r)
}
