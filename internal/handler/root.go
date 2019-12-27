package handler

import (
	"webrtc-server/driver"
	"webrtc-server/internal/handler/http"

	"github.com/gorilla/mux"
)

// RegisterHTTP ...
func RegisterHTTP(db *driver.Database, r *mux.Router) {
	apiPrefix := r.PathPrefix("/api").Subrouter()
	registerAPI(db, apiPrefix)
}

func registerAPI(db *driver.Database, r *mux.Router) {
	userHandler := http.NewUserHandler(db)
	http.RegisterTodo(userHandler, r)
}
