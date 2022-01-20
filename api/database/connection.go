package database

import (
	"fmt"
	"log"

	"github.com/mkrs2404/eKYC/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(dbHost, dbName, dbUser, dbPassword, dbPort string) {

	//Creating DSN string
	connection_string := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s sslmode=disable", dbHost, dbName, dbUser, dbPassword, dbPort)

	//Opening Postgres connection
	connection, err := gorm.Open(postgres.New(postgres.Config{DSN: connection_string}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB = connection

	//Migrating tables to the database
	DB.Debug().AutoMigrate(&models.Plan{}, &models.Client{})

}
