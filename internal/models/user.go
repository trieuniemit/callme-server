package models

// User ...
type User struct {
	Model
	Username string `json:"username" gorm:"UNIQUE"`
	Fullname string `json:"fullname" gorm:"type:text;not null"`
	Password string `json:"password" gorm:"type:text;not null"`
	SocketID string `json:"socket_id"`
	Calling  bool   `json:"calling"`
}
