package http

import (
	"encoding/json"
	"net/http"
	"webrtc-server/driver"
	"webrtc-server/internal/handler/response"
	"webrtc-server/internal/middleware"
	"webrtc-server/internal/models"
	"webrtc-server/internal/repositories"
	"webrtc-server/internal/services"
	"webrtc-server/pkg/helpers"
	"webrtc-server/pkg/jwtauth"

	"github.com/gorilla/mux"
)

type authInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Auth ...
type Auth struct {
	repo       repositories.AuthRepository
	middleware *middleware.Middleware
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
			data := response.Message(true, "success")
			data["user"] = userRegisted

			tokenString, expiresAt, err := jwtauth.CreateToken(user)
			if err != nil {
				data = response.Message(false, err.Error())
				data["key"] = "something_went_wrong"
				response.RespondBadRequest(w, data)
				return
			}

			data = response.Message(true, "success")

			data["token"] = tokenString
			data["expires"] = expiresAt
			data["user"] = user
			response.RespondSuccess(w, data)
			return
		}
		response.RespondSuccess(w, response.Message(false, "Email already exists"))
		return
	}

	response.RespondSuccess(w, response.Message(false, "Register faild!"))
}

// Login user and response token
func (auth *Auth) Login(w http.ResponseWriter, r *http.Request) {
	info := authInfo{}

	err := json.NewDecoder(r.Body).Decode(&info)
	defer r.Body.Close()

	if err == nil {

		userRegisted := auth.repo.Login(info.Email)

		if userRegisted != nil {
			passwordCorrect := helpers.ComparePasswords(userRegisted.Password, info.Password)

			if userRegisted != nil && passwordCorrect {
				data := response.Message(true, "success")
				data["user"] = userRegisted
				data["token"] = "<Token>"
				response.RespondSuccess(w, data)
				return
			}
		}
		response.RespondSuccess(w, response.Message(false, "Email or password incorrect."))
		return
	}

	response.RespondSuccess(w, response.Message(false, "Register faild!"))
}

// NewAuthHandler ...
func NewAuthHandler(db *driver.Database, middleware *middleware.Middleware) *Auth {
	return &Auth{
		repo:       services.NewAuthService(db),
		middleware: middleware,
	}
}

// RegisterAuthRoutes for handle
func RegisterAuthRoutes(a *Auth, routes *mux.Router) {
	routes.HandleFunc("/register", a.Register).Methods("POST")
	routes.HandleFunc("/login", a.Login).Methods("POST")
	//routes.HandleFunc("/logout",  a.middleware.Auth(a.Logout)).Methods("GET")
}
