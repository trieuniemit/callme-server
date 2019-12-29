package models

// User ...
type User struct {
	Model
	Email    string `json:"email" gorm:"UNIQUE"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
}
