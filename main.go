package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mkrs2404/eKYC/api/services"
	"github.com/mkrs2404/eKYC/database"
	"github.com/mkrs2404/eKYC/server"
	"gorm.io/gorm/logger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error fetching the environment values")
	}
	database.Connect(os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), logger.Error)
	services.SeedPlanData()
	server.InitializeRouter()
}
