package models

// User ...
type User struct {
	Model
	Email    int    `json:"email"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
}
