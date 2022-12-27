package entities

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Party struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"-" gorm:"foreignKey:UserID;references:ID"`
	Time        time.Time `json:"datetime"`
	Members     []User    `json:"members,omitempty" gorm:"many2many:party_members;"`
	Location    string    `json:"location"`

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
	Description string `json:"description"`
	// RFC3339 formatted datetime
	Time     time.Time `json:"datetime,omitempty" binding:"required"`
	Location string    `json:"location" binding:"required"`
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

type UpdatePartyRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Location    string    `json:"location" binding:"required"`
	Time        time.Time `json:"datetime" binding:"required"`
}

type JoinInfo struct {
	PartyID uint `json:"party_id"`
	UserID  uint `json:"user_id"`
}
