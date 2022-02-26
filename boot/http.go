package boot

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
	router *gin.Engine
)

func initRoutes() {
	// Ready route
	router.GET("/", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})

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
