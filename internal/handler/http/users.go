package http

import (
	"net/http"
	"webrtc-server/driver"

	"webrtc-server/internal/repositories"
	"webrtc-server/internal/services"

	"github.com/gorilla/mux"
)

// User ...
type User struct {
	repo repositories.UserRepository
}

// NewUserHandler ...
func NewUserHandler(db *driver.Database) *User {
	return &User{
		repo: services.NewUserService(db),
	}
}

// List ...
func (user *User) List(w http.ResponseWriter, r *http.Request) {
	users, err := user.repo.List(10)
	if err != nil {
		data := Message(false, err.Error())
		RespondBadRequest(w, data)
		return
	}
	data := Message(true, "success")
	data["users"] = users

	RespondSuccess(w, data)
}

// Create ..
func (t *User) Create(w http.ResponseWriter, r *http.Request) {
	RespondSuccess(w, Message(true, "success"))
}

// GetByID ..
func (t *User) GetByID(w http.ResponseWriter, r *http.Request) {
	RespondSuccess(w, Message(true, "success"))
}

// Update ..
func (t *User) Update(w http.ResponseWriter, r *http.Request) {
	RespondSuccess(w, Message(true, "success"))
}

// Delete ..
func (t *User) Delete(w http.ResponseWriter, r *http.Request) {
	RespondSuccess(w, Message(true, "success"))
}

// RegisterTodo ..
func RegisterTodo(userHandler *User, routes *mux.Router) {
	routes.HandleFunc("/users", userHandler.List).Methods("GET")
	routes.HandleFunc("/users", userHandler.Create).Methods("POST")
	routes.HandleFunc("/users/{id}", userHandler.GetByID).Methods("GET")
	routes.HandleFunc("/users/{id}", userHandler.Update).Methods("PUT")
	routes.HandleFunc("/users/{id}", userHandler.Delete).Methods("DELETE")
}
