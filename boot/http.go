package boot

import (
	"fmt"
	"github.com/gin-gonic/gin"
	tokens "gitlab.com/ricardo-public/jwt-tools/v2/pkg/token"
	"gitlab.com/ricardo134/party-service/internal/driving/http/party"
	"log"
	"net/http"
)

var (
	router *gin.Engine
)

// @title party-service
// @version 1.0
// @description Ricardo's party service.
//
// @accept json
// @produce json
//
// @contact.name   Ricardo teams
// @contact.email  support@ricardo.net
//
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

func initRoutes() {
	// Ready route
	router.GET("/", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})

	partyController := party.NewController(partyService, []byte(accessSecret))
	tokenMiddleware := tokens.NewJwtAuthMiddleware([]byte(accessSecret))

	partyGroup := router.Group("/parties")
	partyGroup.GET("", tokenMiddleware.Authorize, partyController.Get)
	partyGroup.GET("/user/:user_id", tokenMiddleware.Authorize, partyController.GetForUser)
	partyGroup.GET("/:party_id", tokenMiddleware.Authorize, partyController.GetOne)
	partyGroup.POST("", tokenMiddleware.Authorize, partyController.Create)
	partyGroup.PATCH("/:party_id", tokenMiddleware.Authorize, partyController.Update)
	partyGroup.DELETE("/:party_id", tokenMiddleware.Authorize, partyController.Delete)
}

func ServeHTTP() {
	router = gin.Default()

	initRoutes()

	appURL := fmt.Sprintf("%s:%s", url, port)
	log.Printf("Launching server on %s...\n", appURL)

	log.Fatalln(router.Run(appURL))
}
