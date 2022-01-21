package server

import (
	"github.com/gin-gonic/gin"
)

func InitializeRouter(serverPort string) {

	//Setting up Router
	router := gin.Default()
	InitializeRoutes(router)

	//Starting up Server
	router.Run(serverPort)
}
