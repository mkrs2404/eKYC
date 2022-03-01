package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mkrs2404/eKYC/app/database"
	"github.com/mkrs2404/eKYC/app/helper"
	"github.com/mkrs2404/eKYC/app/minio_client"
	"github.com/mkrs2404/eKYC/app/rabbitmq"
	"github.com/mkrs2404/eKYC/app/redis_client"
	"github.com/mkrs2404/eKYC/app/server"
	_ "github.com/mkrs2404/eKYC/docs"
	"gorm.io/gorm/logger"
)

func init() {
	helper.SetEnvVariablesUtil()
}

// @title           eKYC API
// @version         1.0
// @description     This project provides APIs to sign up a client, upload images, perform face match on face type images or OCR on id card type ones
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  mohit@one2n.in

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
func main() {

	// err := godotenv.Load("../../.env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error fetching the environment values")
	}
	database.Connect(os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), logger.Error)
	// services.SeedPlanData()
	minio_client.InitializeMinio(os.Getenv("MINIO_SERVER"), os.Getenv("MINIO_USER"), os.Getenv("MINIO_PWD"))
	redis_client.InitializeRedis(os.Getenv("REDIS_SERVER"), os.Getenv("REDIS_PASSWORD"))
	rabbitmq.InitializeRabbitMq(os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PWD"), os.Getenv("RABBITMQ_SERVER"))
	server.InitializeRouter()
}
