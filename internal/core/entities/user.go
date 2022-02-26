package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string   `json:"username" gorm:"uniqueIndex, notNull"`
	Parties  []*Party `json:"parties,omitempty" gorm:"many2many:party_members;"`
}
