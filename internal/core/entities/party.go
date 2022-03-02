package entities

import "gorm.io/gorm"

type Party struct {
	gorm.Model
	Name    string  `json:"name,omitempty"`
	UserID  uint    `json:"user_id,omitempty"`
	User    User    `json:"user"`
	Members []*User `json:"members,omitempty" gorm:"many2many:party_members;"`
}

type CreatePartyRequest struct {
	Name   string `json:"name,omitempty" binding:"required"`
	UserID uint   `json:"user_id,omitempty" binding:"required"`
}

type GetPartyRequest struct {
	PartyID uint `json:"party_id,omitempty" binding:"required"`
}

type UpdatePartyRequest struct {
	ID     uint   `json:"id,omitempty" binding:"required"`
	Name   string `json:"name,omitempty" binding:"required"`
	UserID uint   `json:"user_id,omitempty" binding:"required"`
}

type DeletePartyRequest struct {
	ID uint `json:"id,omitempty" binding:"required"`
}
