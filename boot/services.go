package boot

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"ricardo/party-service/internal/core/app"

	"ricardo/party-service/internal/driven/db/cockroachdb"
	"ricardo/party-service/internal/driving/async"
	ricardoNats "ricardo/party-service/internal/driving/async/nats"
)

var (
	partyService app.PartyService
	userService  app.UserService

	natsEncConn  *nats.EncodedConn
	asyncHandler async.Handler
)

func LoadServices() {
	natsConn, err := nats.Connect(fmt.Sprintf("nats://%s:%s@%s", natsUsr, natsPwd, natsURL))
	if err != nil {
		log.Fatal(err)
	}
	natsEncConn, err = nats.NewEncodedConn(natsConn, nats.JSON_ENCODER)

	partyRepo := cockroachdb.NewPartyRepository(client)
	partyService = app.NewPartyService(partyRepo)

	userRepo := cockroachdb.NewUserRepository(client)
	userService = app.NewUserService(userRepo)

	asyncHandler = ricardoNats.NewNatsUserHandler(partyService, userService)
}
