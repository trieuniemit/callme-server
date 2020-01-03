package models

// CallHistory struct model
type CallHistory struct {
	Model
	Length  int  `json:"length"`
	Missing bool `json:"missing" gorm:"default:true"`
}
