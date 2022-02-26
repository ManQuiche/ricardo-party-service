package entities

import "gorm.io/gorm"

type Party struct {
	gorm.Model `json:"gorm_._model"`
	Name       string  `json:"name,omitempty"`
	UserID     uint    `json:"user_id,omitempty"`
	User       User    `json:"user"`
	Members    []*User `json:"members,omitempty" gorm:"many2many:party_members;"`
}
