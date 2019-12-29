package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"webrtc-server/driver"
	"webrtc-server/pkg/helpers"

	"webrtc-server/internal/models"
	"webrtc-server/internal/repositories"
	"webrtc-server/internal/services"

	"github.com/gorilla/mux"
)

// User ...
type User struct {
	repo repositories.UserRepository
}

// List ...
func (u *User) List(w http.ResponseWriter, r *http.Request) {
	users, err := u.repo.List(10)
	if err != nil {
		data := Message(false, err.Error())
		RespondBadRequest(w, data)
		return
	}
	data := Message(true, "success")
	data["users"] = users

	RespondSuccess(w, data)
}

// Create new user
func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()

	if err == nil {
		passwordHash, _ := helpers.HashAndSalt(user.Password)
		user.Password = passwordHash
		u.repo.Create(&user)

		data := Message(true, "success")
		data["user"] = user
		RespondSuccess(w, data)
		return
	}

	data := Message(false, "Field is required")
	RespondBadRequest(w, data)
}

// GetByID return user
func (u *User) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if i, err := strconv.ParseInt(vars["id"], 10, 64); err == nil {
		user, err := u.repo.GetByID(i)

		if err != nil {
			data := Message(false, err.Error())
			RespondBadRequest(w, data)
			return
		}

		data := Message(true, "success")
		data["user"] = user
		RespondSuccess(w, data)
	}
}

// Update user
func (u *User) Update(w http.ResponseWriter, r *http.Request) {

	user := models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()

	if err == nil {
		vars := mux.Vars(r)

		if id, err := strconv.ParseInt(vars["id"], 10, 64); err == nil {
			userFound, err := u.repo.GetByID(id)

			if err != nil {
				data := Message(false, err.Error())
				RespondBadRequest(w, data)
				return
			}

			passwordHash, err := helpers.HashAndSalt(user.Password)
			userFound.Password = passwordHash

			userFound.Email = user.Email
			userFound.Fullname = user.Fullname
			u.repo.Update(userFound)

			data := Message(true, "success")
			data["user"] = userFound
			RespondSuccess(w, data)
			return
		}
	}

	data := Message(false, "Field is required")
	RespondBadRequest(w, data)
}

// Delete an user
func (u *User) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if id, err := strconv.Atoi(vars["id"]); err == nil {
		status, _ := u.repo.Delete(int64(id))
		if status == true {
			RespondSuccess(w, Message(true, "success"))
			return
		}
	}

	RespondBadRequest(w, Message(true, "User not found."))
}

// NewUserHandler ...
func NewUserHandler(db *driver.Database) *User {
	return &User{
		repo: services.NewUserService(db),
	}
}

// RegisterUserRoutes for handle
func RegisterUserRoutes(userHandler *User, routes *mux.Router) {
	routes.HandleFunc("/users", userHandler.List).Methods("GET")
	routes.HandleFunc("/users", userHandler.Create).Methods("POST")
	routes.HandleFunc("/users/{id}", userHandler.GetByID).Methods("GET")
	routes.HandleFunc("/users/{id}", userHandler.Update).Methods("PUT")
	routes.HandleFunc("/users/{id}", userHandler.Delete).Methods("DELETE")
}
