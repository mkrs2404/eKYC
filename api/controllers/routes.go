package controllers

import "github.com/gin-gonic/gin"

func InitializeRoutes(router *gin.Engine) {

	//Signup API routes
	router.POST("/api/v1/signup", signUp)

}
