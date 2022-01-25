package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/api/controllers"
	"github.com/mkrs2404/eKYC/api/middlewares"
)

func InitializeRoutes(router *gin.Engine) {

	//Default routes
	router.GET("/", controllers.WelcomePage)
	router.NoRoute(controllers.NoRoute)

	//Signup API routes
	router.POST("/api/v1/signup", controllers.SignUpClient)

	//Image Upload API routes
	router.POST("/api/v1/image", middlewares.AuthRequired(), controllers.UploadImageClient)

}
