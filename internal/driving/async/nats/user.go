package nats

import (
	"context"
	"ricardo/party-service/internal/core/app/party"
	"ricardo/party-service/internal/driving/async"
)

type natsHandler struct {
	partyService party.Service
}

func NewNatsUserHandler(svc party.Service) async.Handler {
	return natsHandler{svc}
}

func (nh natsHandler) OnUserDelete(userID uint) {
	_ = nh.partyService.DeleteAllForUser(context.Background(), userID)
}
