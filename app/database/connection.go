package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

//Connect connects to the database and assigns the connection to global variable DB
func Connect(dbHost string, dbName string, dbUser string, dbPassword string, dbPort string, loggerLevel logger.LogLevel) {

	//Creating DSN string
	connection_string := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s sslmode=disable", dbHost, dbName, dbUser, dbPassword, dbPort)

	//Opening Postgres connection
	connection, err := gorm.Open(postgres.New(postgres.Config{DSN: connection_string}), &gorm.Config{
		Logger: logger.Default.LogMode(loggerLevel),
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
}
