package models

// User ...
type User struct {
	Model
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
}
