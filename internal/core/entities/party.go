package entities

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Party struct {
	gorm.Model
	Name    string    `json:"name,omitempty"`
	UserID  uint      `json:"user_id,omitempty"`
	User    User      `json:"-"`
	Time    time.Time `json:"time,omitempty"`
	Members []*User   `json:"members,omitempty" gorm:"many2many:party_members;"`
}

func (p *Party) MarshalJSON() ([]byte, error) {
	type Alias Party
	return json.Marshal(&struct {
		*Alias
		Time string `json:"time,omitempty"`
	}{
		Alias: (*Alias)(p),
		Time:  p.Time.Format(time.RFC3339),
	})
}

type CreatePartyRequest struct {
	Name   string `json:"name,omitempty" binding:"required"`
	UserID uint   `json:"user_id,omitempty" binding:"required"`
	// RFC3339 formatted datetime
	Time time.Time `json:"time,omitempty" binding:"required"`
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
	ID     uint   `json:"id,omitempty" binding:"required"`
	Name   string `json:"name,omitempty" binding:"required"`
	UserID uint   `json:"user_id,omitempty" binding:"required"`
}

type DeletePartyRequest struct {
	ID uint `json:"id,omitempty" binding:"required"`
}
