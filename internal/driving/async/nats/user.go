package nats

import (
	"context"
	"gitlab.com/ricardo134/party-service/internal/core/app"
	"gitlab.com/ricardo134/party-service/internal/core/entities"
)

type handler struct {
	partyService app.PartyService
	userService  app.UserService
}

type Handler interface {
	Created(user entities.User)
	Updated(user entities.User)
	Deleted(userID uint)
}

func NewUserHandler(partySvc app.PartyService, userSvc app.UserService) Handler {
	return handler{partySvc, userSvc}
}

func (nh handler) Created(user entities.User) {
	_, _ = nh.userService.Save(context.Background(), user)
}

func (nh handler) Updated(user entities.User) {
	_, _ = nh.userService.Save(context.Background(), user)
}

func (nh handler) Deleted(userID uint) {
	_ = nh.partyService.DeleteAllForUser(context.Background(), userID)
	_ = nh.userService.Delete(context.Background(), userID)
}
