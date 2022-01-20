package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mkrs2404/eKYC/api/database"
	"github.com/mkrs2404/eKYC/api/server"
)

func main() {
	//api.Run()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error fetching the environment values")
	} else {
		database.Connect(os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"))
		server.InitializeRouter(os.Getenv("SERVER_ADDR"))
	}
}
