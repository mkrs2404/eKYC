package database

import (
	"fmt"
	"log"
	"time"

	"github.com/mkrs2404/eKYC/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

//Connect connects to the database and assigns the connection to global variable DB
func Connect(dbHost string, dbName string, dbUser string, dbPassword string, dbPort string) {

	//Creating DSN string
	connection_string := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s sslmode=disable", dbHost, dbName, dbUser, dbPassword, dbPort)

	//Opening Postgres connection
	connection, err := gorm.Open(postgres.New(postgres.Config{DSN: connection_string}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	db, _ := connection.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.SetConnMaxLifetime(time.Hour / 2)

	if err != nil {
		log.Fatal(err)
	}
	DB = connection

	//Migrating tables to the database
	DB.Debug().AutoMigrate(&models.Plan{}, &models.Client{})

}
