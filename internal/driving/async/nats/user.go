package nats

import (
	"context"
	"gorm.io/gorm"
	"ricardo/party-service/internal/core/app"
	"ricardo/party-service/internal/core/entities"
	"ricardo/party-service/internal/driving/async"
)

type natsHandler struct {
	partyService app.PartyService
	userService  app.UserService
}

func NewNatsUserHandler(partySvc app.PartyService, userSvc app.UserService) async.Handler {
	return natsHandler{partySvc, userSvc}
}

func (nh natsHandler) OnUserDelete(userID uint) {
	_ = nh.partyService.DeleteAllForUser(context.Background(), userID)
	_ = nh.userService.Delete(context.Background(), userID)
}

func (nh natsHandler) OnAccountCreated(userID uint) {
	_, _ = nh.userService.Save(context.Background(), entities.User{
		Model: gorm.Model{
			ID: userID,
		},
		Parties: nil,
	})
}
