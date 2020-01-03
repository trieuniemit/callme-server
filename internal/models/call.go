package models

// Call struct model
type Call struct {
	Model
	UserID   uint `json:"user_id" `
	TargetID uint `json:"to_id" `
	User     User `gorm:"foreignkey:UserID"`
	Target   User `json:"-" gorm:"foreignkey:TargetID"`
}
