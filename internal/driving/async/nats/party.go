package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/nats-io/nats.go"
	"gitlab.com/ricardo134/party-service/internal/core/app"
	"gitlab.com/ricardo134/party-service/internal/core/entities"
)

type partyHandler struct {
	partyService app.PartyService
}

type PartyHandler interface {
	Joined(msg *nats.Msg)
	Requested(msg *nats.Msg)
}

func NewPartyHandler(partySvc app.PartyService) PartyHandler {
	return partyHandler{partySvc}
}

func marshalErr(err error) []byte {
	if err != nil {
		jErr, marshalErr := json.Marshal(err)
		if marshalErr != nil {
			err2 := fmt.Errorf("cannot marshal error: %w", err)
			log.Println(err2)
			jErr, _ = json.Marshal(err2)
		}

		return jErr
	}

	return nil
}

func respondErr(msg *nats.Msg, err error) {
	jErr := marshalErr(err)

	// No response if nil so it times out publisher side
	if jErr != nil {
		_ = msg.Respond(jErr)
	}
}

func (p partyHandler) Joined(msg *nats.Msg) {
	var info entities.JoinInfo
	err := json.Unmarshal(msg.Data, &info)
	if err != nil {
		err = fmt.Errorf("unmarshalling join info: %w", err)
		log.Println(err)
		respondErr(msg, err)
		return
	}

	err = p.partyService.Joined(context.Background(), info.PartyID, info.UserID)
	if err != nil {
		err = fmt.Errorf("joining party: %w", err)
		log.Println(err)
		respondErr(msg, err)
	}
}

func (p partyHandler) Requested(msg *nats.Msg) {
	var partyID int
	partyID, err := strconv.Atoi(string(msg.Data))
	if err != nil {
		err = fmt.Errorf("unmarshalling nats data: %w", err)
		log.Println(err)
		respondErr(msg, err)
		return
	}

	party, err := p.partyService.Get(context.Background(), uint(partyID))
	if err != nil {
		err = fmt.Errorf("fetching requested party: %w", err)
		log.Println(err)
		respondErr(msg, err)
		return
	}

	jsonParty, err := json.Marshal(party)
	if err != nil {
		err = fmt.Errorf("marshalling requested party: %w", err)
		log.Println(err)
		respondErr(msg, err)
		return
	}

	_ = msg.Respond(jsonParty)
}
