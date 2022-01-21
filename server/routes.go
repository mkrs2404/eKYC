package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/api/controllers"
)

func InitializeRoutes(router *gin.Engine) {

	//Signup API routes
	router.POST("/api/v1/signup", controllers.SignUp)

}
