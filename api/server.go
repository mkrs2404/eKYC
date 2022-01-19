package api

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mkrs2404/eKYC/api/controllers"
)

var server = controllers.Server{}

func Run() {
	//Loading env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error fetching the environment values")
	} else {
		server.InitializeDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"))
		server.InitializeRouter(os.Getenv("SERVER_ADDR"))
	}
}
