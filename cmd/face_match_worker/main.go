package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mkrs2404/eKYC/app/helper"
	"github.com/mkrs2404/eKYC/app/rabbitmq"
	"github.com/mkrs2404/eKYC/app/redis_client"
	"github.com/mkrs2404/eKYC/app/workers"
)

func init() {
	helper.SetEnvVariablesUtil()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error fetching the environment values")
	}
	redis_client.InitializeRedis(os.Getenv("REDIS_SERVER"), os.Getenv("REDIS_PASSWORD"))
	rabbitmq.InitializeRabbitMq(os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PWD"), os.Getenv("RABBITMQ_SERVER"))
	workers.FaceMatchUtil()
}
