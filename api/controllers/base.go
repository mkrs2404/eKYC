package controllers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func (server *Server) InitializeDatabase(dbHost, dbName, dbUser, dbPassword, dbPort string) {
	var err error

	//Creating DSN string
	connection_string := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s sslmode=disable", dbHost, dbName, dbUser, dbPassword, dbPort)

	//Opening Postgres connection
	server.DB, err = gorm.Open(postgres.New(postgres.Config{DSN: connection_string}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	//Migrating tables to the database
	server.DB.Debug().AutoMigrate(&models.Plan{}, &models.Client{})

}

func (server *Server) InitializeRouter(serverPort string) {

	//Setting up Router
	server.Router = gin.Default()
	server.InitializeRoutes()

	//Starting up Server
	server.Router.Run(serverPort)
}
