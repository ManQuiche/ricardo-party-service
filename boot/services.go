package boot

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"gitlab.com/ricardo134/party-service/internal/core/app"
	"log"

	"gitlab.com/ricardo134/party-service/internal/driven/db/postgresql"
	natsext "gitlab.com/ricardo134/party-service/internal/driving/async/nats"
)

var (
	partyService app.PartyService
	userService  app.UserService

	natsEncConn *nats.EncodedConn
	userAsync   natsext.UserHandler
	partyAsync  natsext.PartyHandler
)

func LoadServices() {
	natsConn, err := nats.Connect(fmt.Sprintf("nats://%s:%s@%s", natsUsr, natsPwd, natsURL))
	if err != nil {
		log.Fatal(err)
	}
	natsEncConn, err = nats.NewEncodedConn(natsConn, nats.JSON_ENCODER)

	partyRepo := postgresql.NewPartyRepository(client)
	partyService = app.NewPartyService(partyRepo)

	userRepo := postgresql.NewUserRepository(client)
	userService = app.NewUserService(userRepo)

	userAsync = natsext.NewUserHandler(partyService, userService)
	partyAsync = natsext.NewPartyHandler(partyService)
}
