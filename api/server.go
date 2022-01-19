package api

import (
	"fmt"
	"log"

	"github.com/mkrs2404/eKYC/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB *gorm.DB
}

func (server *Server) InitializeDatabase(DbHost, DbName, DbUser, DbPassword, DbPort string) {
	var err error

	//Creating DSN string
	connection_string := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s sslmode=disable", DbHost, DbName, DbUser, DbPassword, DbPort)

	//Opening Postgres connection
	server.DB, err = gorm.Open(postgres.New(postgres.Config{DSN: connection_string}), &gorm.Config{})
	HandleError(err)

	//Migrating tables to the database
	server.DB.Debug().AutoMigrate(&models.Plan{}, &models.Client{})

}

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
