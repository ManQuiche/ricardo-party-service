package boot

import (
	"ricardo/party-service/internal/core/app/party"
	"ricardo/party-service/internal/driven/db/cockroachdb"
)

var (
	partyService party.Service
)

func LoadServices() {
	partyRepo := cockroachdb.NewPartyRepository(client)
	partyService = party.NewPartyService(partyRepo)
}
