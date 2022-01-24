package server

import (
	"os"

	"github.com/gin-gonic/gin"
)

func InitializeRouter() {
	serverAddr := os.Getenv("SERVER_ADDR")

	//Setting up Router
	router := gin.Default()
	InitializeRoutes(router)

	//Starting up Server
	router.Run(serverAddr)
}
