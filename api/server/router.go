package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/api/controllers"
)

func InitializeRouter(serverPort string) {

	//Setting up Router
	router := gin.Default()
	controllers.InitializeRoutes(router)

	//Starting up Server
	router.Run(serverPort)
}
