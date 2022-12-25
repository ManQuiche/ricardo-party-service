package nats

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"gitlab.com/ricardo134/party-service/internal/core/app"
	"log"
	"strconv"
)

type partyHandler struct {
	partyService app.PartyService
}

type PartyHandler interface {
	Joined(partyID, userID uint)
	Requested(msg *nats.Msg)
}

func NewPartyHandler(partySvc app.PartyService) PartyHandler {
	return partyHandler{partySvc}
}

func (p partyHandler) Joined(partyID, userID uint) {
	err := p.partyService.Joined(context.Background(), partyID, userID)
	if err != nil {
		log.Print(err.Error())
	}
}

func (p partyHandler) Requested(msg *nats.Msg) {
	var partyID int
	partyID, err := strconv.Atoi(string(msg.Data))
	if err != nil {
		_ = msg.Respond(nil)
		return
	}

	party, err := p.partyService.Get(context.Background(), uint(partyID))
	if err != nil {
		_ = msg.Respond(nil)
		return
	}

	jsonParty, err := json.Marshal(party)
	if err != nil {
		log.Print(err.Error())
		_ = msg.Respond(nil)
		return
	}

	_ = msg.Respond(jsonParty)
}
