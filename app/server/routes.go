package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/app/controllers"
	"github.com/mkrs2404/eKYC/app/middlewares"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func InitializeRoutes(router *gin.Engine) {

	//Default routes
	router.GET("/", controllers.WelcomePage)
	router.NoRoute(controllers.NoRoute)

	//Signup API routes
	router.POST("/api/v1/signup", controllers.SignUpClient)

	//Async APIs
	router.POST("/api/v1/face-match-async", middlewares.AuthRequired(), controllers.AsyncFaceMatchClient, middlewares.UpdateApi())
	router.POST("/api/v1/get-score", middlewares.AuthRequired(), controllers.GetFaceMatchScore, middlewares.UpdateApi())

	router.POST("/api/v1/ocr-async", middlewares.AuthRequired(), controllers.AsyncOcrClient, middlewares.UpdateApi())
	router.POST("/api/v1/get-ocr-data", middlewares.AuthRequired(), controllers.GetOcrData, middlewares.UpdateApi())

	routerGroup := router.Group("/api/v1")
	authRouterGroup := routerGroup.Use(middlewares.AuthRequired())
	ImageAPI(authRouterGroup)
	MatchAPI(authRouterGroup)
	OcrAPI(authRouterGroup)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func ImageAPI(r gin.IRoutes) {
	r.POST("/image", controllers.UploadImageClient)
}

func MatchAPI(r gin.IRoutes) {
	r.Use(middlewares.SaveApi())
	r.POST("/face-match", controllers.FaceMatchClient)
}

func OcrAPI(r gin.IRoutes) {
	r.Use(middlewares.SaveApi())
	r.POST("/ocr", controllers.OcrClient)
}
