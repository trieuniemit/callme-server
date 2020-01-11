package models

// Contact struct
type Contact struct {
	ID    uint `gorm:"primary_key"`
	User1 uint
	User2 uint
}
