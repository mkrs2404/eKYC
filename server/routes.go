package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/api/controllers"
)

func InitializeRoutes(router *gin.Engine) {

	//Default routes
	router.POST("/", controllers.WelcomePage)
	router.GET("/", controllers.WelcomePage)
	router.NoRoute(controllers.NoRoute)

	//Signup API routes
	router.POST("/api/v1/signup", controllers.SignUp)

}
