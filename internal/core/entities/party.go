package entities

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Party struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	UserID      uint      `json:"user_id,omitempty"`
	User        User      `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Time        time.Time `json:"datetime,omitempty"`
	Members     []User    `json:"members,omitempty" gorm:"many2many:party_members;"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (p *Party) MarshalJSON() ([]byte, error) {
	type Alias Party
	return json.Marshal(&struct {
		*Alias
		Time string `json:"datetime,omitempty"`
	}{
		Alias: (*Alias)(p),
		Time:  p.Time.Format(time.RFC3339),
	})
}

type CreatePartyRequest struct {
	Name        string `json:"name,omitempty" binding:"required"`
	Description string `json:"description,omitempty"`
	// RFC3339 formatted datetime
	Time time.Time `json:"datetime,omitempty" binding:"required"`
}

// TODO: inspect
//func (cpr *CreatePartyRequest) UnmarshalJSON(data []byte) error {
//	type cprJSON struct {
//		Name   string `json:"name,omitempty"`
//		UserID uint   `json:"user_id,omitempty"`
//		Time   string `json:"time,omitempty"`
//	}
//	type Alias CreatePartyRequest
//	var cprj Alias
//	if err := json.Unmarshal(data, &cprj); err != nil {
//		return err
//	}
//
//	cpr.Name = cprj.Name
//	cpr.UserID = cprj.UserID
//	sTime := cprj.Time.Format(time.RFC3339)
//	print(sTime)
//	//pTime, err := time.Parse(time.RFC3339, cprj.Time)
//	//if err != nil {
//	//	return err
//	//}
//	//cpr.Time = pTime
//
//	return nil
//}

type GetPartyRequest struct {
	PartyID uint `json:"party_id,omitempty"`
	UserID  uint `json:"user_id,omitempty"`
}

type UpdatePartyRequest struct {
	Name   string `json:"name,omitempty" binding:"required"`
	UserID uint   `json:"user_id,omitempty" binding:"required"`
}
