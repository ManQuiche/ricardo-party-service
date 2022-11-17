package nats

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/nats-io/nats.go"
	"gitlab.com/ricardo134/party-service/internal/core/app"
	"log"
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
	//TODO implement me
	panic("implement me")
}

func (p partyHandler) Requested(msg *nats.Msg) {
	var partyID uint
	buf := bytes.NewReader(msg.Data)
	err := binary.Read(buf, binary.BigEndian, &partyID)
	if err != nil {
		_ = msg.Respond(nil)
	}

	party, err := p.partyService.Get(context.Background(), partyID)
	if err != nil {
		_ = msg.Respond(nil)
	}

	resbuf := new(bytes.Buffer)
	err = binary.Write(resbuf, binary.BigEndian, party)
	if err != nil {
		log.Print(err.Error())
		_ = msg.Respond(nil)
	}

	_ = msg.Respond(resbuf.Bytes())
}
