package models

// User ...
type User struct {
	Model
	Email    string `json:"email" gorm:"UNIQUE"`
	Fullname string `json:"fullname" gorm:"type:text;not null"`
	Password string `json:"-" gorm:"type:text;not null"`
	SocketID string `json:"socket_id"`
	Calling  bool   `json:"calling"`
}
