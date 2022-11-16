package nats

import "gitlab.com/ricardo134/party-service/internal/core/app"

type partyHandler struct {
	partyService app.PartyService
}

type PartyHandler interface {
	Joined(partyID uint)
	Requested(partyID, userID uint)
}

func NewPartyHandler(partySvc app.PartyService) PartyHandler {
	return partyHandler{partySvc}
}
