package boot

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

	partyGroup := router.Group("/parties")
	partyGroup.GET("", partyController.Get)
	partyGroup.GET("/user/:user_id", partyController.GetForUser)
	partyGroup.GET("/:party_id", partyController.GetOne)
	partyGroup.POST("", partyController.Create)
	partyGroup.PATCH("", partyController.Update)
	partyGroup.DELETE("", partyController.Delete)
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
