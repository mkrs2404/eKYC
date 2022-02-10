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
	"github.com/mkrs2404/eKYC/app/services"
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
	redis_addr     = kingpin.Flag("redis_server", "Redis Server_IP:Port").String()
	redis_pwd      = kingpin.Flag("redis_pwd", "Redis Password").String()
)

func init() {
	kingpin.Parse()
	helper.SetEnvVariablesUtil(*hostname, *dbname, *user, *password, *port, *server_address, *minio_server, *minio_user, *minio_pwd, *redis_addr, *redis_pwd)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error fetching the environment values")
	}
	database.Connect(os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), logger.Error)
	services.SeedPlanData()
	minio_client.InitializeMinio(os.Getenv("MINIO_SERVER"), os.Getenv("MINIO_USER"), os.Getenv("MINIO_PWD"))
	redis_client.InitializeRedis(os.Getenv("REDIS_SERVER"), os.Getenv("REDIS_PASSWORD"))
	server.InitializeRouter()
}
