package http

import (
	"encoding/json"
	"net/http"
	"webrtc-server/driver"
	"webrtc-server/internal/models"
	"webrtc-server/internal/repositories"
	"webrtc-server/internal/services"
	"webrtc-server/pkg/helpers"

	"github.com/gorilla/mux"
)

type authInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Auth ...
type Auth struct {
	repo repositories.AuthRepository
}

// Register new account
func (auth *Auth) Register(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()

	if err == nil {
		passwordHash, _ := helpers.HashAndSalt(user.Password)
		user.Password = passwordHash
		userRegisted := auth.repo.Register(&user)
		if userRegisted != nil {
			data := Message(true, "success")
			data["user"] = userRegisted
			RespondSuccess(w, data)
			return
		}
		RespondSuccess(w, Message(false, "Email already exists"))
		return
	}

	RespondSuccess(w, Message(false, "Register faild!"))
}

// Login ..
func (auth *Auth) Login(w http.ResponseWriter, r *http.Request) {
	info := authInfo{}

	err := json.NewDecoder(r.Body).Decode(&info)
	defer r.Body.Close()

	if err == nil {
		passwordHash, _ := helpers.HashAndSalt(info.Password)
		info.Password = passwordHash

		userRegisted := auth.repo.Login(info.Email)

		if userRegisted != nil {
			data := Message(true, "success")
			data["user"] = userRegisted
			RespondSuccess(w, data)
			return
		}
		RespondSuccess(w, Message(false, "Email already exists"))
		return
	}

	RespondSuccess(w, Message(false, "Register faild!"))
}

// NewAuthHandler ...
func NewAuthHandler(db *driver.Database) *Auth {
	return &Auth{
		repo: services.NewAuthService(db),
	}
}

// RegisterAuthRoutes for handle
func RegisterAuthRoutes(authHandler *Auth, routes *mux.Router) {
	routes.HandleFunc("/register", authHandler.Register).Methods("POST")
	// routes.HandleFunc("/login", authHandler.Login).Methods("POST")
	// routes.HandleFunc("/logout", authHandler.Logout).Methods("GET")
}
