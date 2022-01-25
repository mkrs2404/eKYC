package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mkrs2404/eKYC/api/services"
	"github.com/mkrs2404/eKYC/database"
	_ "github.com/mkrs2404/eKYC/minio_client"
	"github.com/mkrs2404/eKYC/server"
	"gopkg.in/alecthomas/kingpin.v2"
	"gorm.io/gorm/logger"
)

func init() {
	var (
		hostname       = kingpin.Flag("host", "Hostname").Required().String()
		dbname         = kingpin.Flag("db", "Database name").Required().String()
		user           = kingpin.Flag("user", "Username").Required().String()
		password       = kingpin.Flag("pwd", "Password").Required().String()
		port           = kingpin.Flag("port", "Database port").Required().String()
		server_address = kingpin.Flag("server", "Server_IP:Port").Required().String()
	)
	kingpin.Parse()
	os.Setenv("DB_HOST", *hostname)
	os.Setenv("DB_NAME", *dbname)
	os.Setenv("DB_USER", *user)
	os.Setenv("DB_PASSWORD", *password)
	os.Setenv("DB_PORT", *port)
	os.Setenv("SERVER_ADDR", *server_address)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error fetching the environment values")
	}
	database.Connect(os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), logger.Error)
	services.SeedPlanData()
	server.InitializeRouter()
}
