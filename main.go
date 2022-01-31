package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mkrs2404/eKYC/api/services"
	"github.com/mkrs2404/eKYC/database"
	"github.com/mkrs2404/eKYC/minio_client"
	"github.com/mkrs2404/eKYC/server"
	"gopkg.in/alecthomas/kingpin.v2"
	"gorm.io/gorm/logger"
)

var (
	hostname       = kingpin.Flag("host", "Hostname").String()
	dbname         = kingpin.Flag("db", "Database name").String()
	user           = kingpin.Flag("user", "Username").String()
	password       = kingpin.Flag("pwd", "Password").String()
	port           = kingpin.Flag("port", "Database port").String()
	server_address = kingpin.Flag("server", "Server_IP:Port").String()
	minio_server   = kingpin.Flag("minio_server", "Minio Server_IP:Port").String()
	minio_pwd      = kingpin.Flag("minio_pwd", "Minio Password").String()
	minio_user     = kingpin.Flag("minio_user", "Minio user").String()
)

func init() {
	kingpin.Parse()
	SetEnvVariablesUtil()
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error fetching the environment values")
	}
	database.Connect(os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), logger.Error)
	services.SeedPlanData()
	server.InitializeRouter()
	minio_client.InitializeMinio(os.Getenv("MINIO_SERVER"), os.Getenv("MINIO_USER"), os.Getenv("MINIO_PWD"))
}
