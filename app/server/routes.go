package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/app/controllers"
	"github.com/mkrs2404/eKYC/app/middlewares"
)

func InitializeRoutes(router *gin.Engine) {

	//Default routes
	router.GET("/", controllers.WelcomePage)
	router.NoRoute(controllers.NoRoute)

	//Signup API routes
	router.POST("/api/v1/signup", controllers.SignUpClient)

	routerGroup := router.Group("/api/v1")
	authRouterGroup := routerGroup.Use(middlewares.AuthRequired())
	ImageAPI(authRouterGroup)
	MatchAPI(authRouterGroup)
	OcrAPI(authRouterGroup)
}

func ImageAPI(r gin.IRoutes) {
	r.POST("/image", controllers.UploadImageClient)
}

func MatchAPI(r gin.IRoutes) {
	r.Use(middlewares.SaveApi())
	r.POST("/face-match", controllers.FaceMatchClient)
	r.POST("/face-match-async", controllers.AsyncFaceMatchClient)
	r.POST("/get-score", controllers.GetFaceMatchScore)

}

func OcrAPI(r gin.IRoutes) {
	r.Use(middlewares.SaveApi())
	r.POST("/ocr", controllers.OcrClient)
	r.POST("/ocr-async", controllers.AsyncOcrClient)
	r.POST("/get-ocr-data", controllers.GetOcrData)
}
