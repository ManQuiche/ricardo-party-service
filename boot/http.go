package boot

import (
	"fmt"
	"github.com/gin-gonic/gin"
	tokens "gitlab.com/ricardo-public/jwt-tools/pkg"
	"log"
	"net/http"
	"ricardo/party-service/internal/driving/http/party"
)

var (
	router *gin.Engine
)

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
	partyGroup.PATCH("", tokenMiddleware.Authorize, partyController.Update)
	partyGroup.DELETE("", tokenMiddleware.Authorize, partyController.Delete)
}

func ServeHTTP() {
	router = gin.Default()

	initRoutes()

	appURL := fmt.Sprintf("%s:%s", url, port)
	log.Printf("Launching server on %s...\n", appURL)

	log.Fatalln(router.Run(appURL))

	// TODO: go func and etc
	//log.Println("HTTP server stopped, exiting...")
}
