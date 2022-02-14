package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mkrs2404/eKYC/app/database"
	"github.com/mkrs2404/eKYC/app/helper"
	"github.com/mkrs2404/eKYC/app/minio_client"
	"github.com/mkrs2404/eKYC/app/redis_client"
	"github.com/mkrs2404/eKYC/app/server"
	"gorm.io/gorm/logger"
)

func init() {
	helper.SetEnvVariablesUtil()
}

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
	server.InitializeRouter()
}
