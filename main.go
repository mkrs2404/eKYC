package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	connection_string := "host=localhost user=postgres password=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: connection_string}), &gorm.Config{})
	handleError(err)

	database, err := db.DB()
	handleError(err)

	err = database.Ping()
	handleError(err)

	fmt.Println("Connected to Postgres")

}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
