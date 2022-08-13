package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Parties []Party `json:"parties,omitempty" gorm:"many2many:party_members;"`
	Test    int     `json:"test,omitempty"`
}
