package handler

import (
	"webrtc-server/driver"
	"webrtc-server/internal/handler/http"
	"webrtc-server/internal/handler/socket"
	"webrtc-server/internal/middleware"

	"github.com/gorilla/mux"
)

// RegisterHTTP ...
func RegisterHTTP(db *driver.Database, r *mux.Router) {
	apiPrefix := r.PathPrefix("/api").Subrouter()

	middleware := middleware.NewMiddleware(db)

	registerAPI(db, apiPrefix, middleware)

	// register websocket
	socket.InitSocketRoute(socket.NewSocketHandler(db), r)
}

func registerAPI(db *driver.Database, r *mux.Router, middleware *middleware.Middleware) {
	// register api
	http.RegisterAuthRoutes(http.NewAuthHandler(db, middleware), r)
	http.RegisterUserRoutes(http.NewUserHandler(db, middleware), r)
}
