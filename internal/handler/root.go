package handler

import (
	"tracker/driver"
	"tracker/internal/handler/http"

	"github.com/gorilla/mux"
)

func RegisterHTTP(db *driver.Database, r *mux.Router) {
	apiPrefix := r.PathPrefix("/api").Subrouter()
	registerAPI(db, apiPrefix)
}

func registerAPI(db *driver.Database, r *mux.Router) {
	todoHandler := http.NewTodoHandler(db)

	http.RegisterTodo(todoHandler, r)
}
