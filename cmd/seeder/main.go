package main

import (
	"log"
	"os"

	"github.com/danvergara/seeder"
	"github.com/joho/godotenv"
	"github.com/mkrs2404/eKYC/app/database"
	"github.com/mkrs2404/eKYC/app/database/seeds"
	"github.com/mkrs2404/eKYC/app/helper"
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
	database.Connect(os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), logger.Error)
	s := seeds.NewSeed(database.DB)

	if err := seeder.Execute(s); err != nil {
		log.Fatalf("error seeding the db %s\n", err)
	}
}
