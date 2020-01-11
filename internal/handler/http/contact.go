package http

import (
	"net/http"
	"webrtc-server/driver"
	"webrtc-server/internal/handler/response"
	"webrtc-server/internal/middleware"
	"webrtc-server/internal/models"
	"webrtc-server/internal/repositories"
	"webrtc-server/internal/services"

	"github.com/gorilla/mux"
)

// Contact ...
type Contact struct {
	repo repositories.ContactRepository
}

// GetList new account
func (c *Contact) GetList(w http.ResponseWriter, r *http.Request) {

	currentUser := r.Context().Value("user").(models.User)

	users, err := c.repo.GetList(&currentUser)

	if err != nil {
		response.RespondSuccess(w, response.Message(false, "Faild to get list contact!"))
		return
	}

	data := response.Message(true, "success")
	data["users"] = users
	response.RespondSuccess(w, data)

}

// AddNew ...
func (c *Contact) AddNew(w http.ResponseWriter, r *http.Request) {
	//info := authInfo{}

	// err := json.NewDecoder(r.Body).Decode(&info)
	// defer r.Body.Close()

	// response.RespondBadRequest(w, response.Message(false, "Register faild!"))
}

// NewContactHandler ...
func NewContactHandler(db *driver.Database) *Contact {
	return &Contact{
		repo: services.NewContactService(db),
	}
}

// RegisterContactRoutes for handle
func RegisterContactRoutes(c *Contact, routes *mux.Router, db *driver.Database) {
	routes.HandleFunc("/contact", middleware.Authenticate(c.GetList, db)).Methods("GET")
	routes.HandleFunc("/contact", middleware.Authenticate(c.AddNew, db)).Methods("POST")
	//routes.HandleFunc("/logout",  a.middleware.Auth(a.Logout)).Methods("GET")
}
