package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mkrs2404/eKYC/app/database"
	"github.com/mkrs2404/eKYC/app/helper"
	"github.com/mkrs2404/eKYC/app/models"
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

	//Migrating tables to the database
	database.DB.AutoMigrate(&models.Plan{}, &models.Client{}, &models.File{}, &models.Api_Calls{})

}
