package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var rabbitMq *amqp.Connection

func InitializeRabbitMq(userName, password, host string) {
	if userName == "" || password == "" {
		userName, password = "guest", "guest"
	}
	connString := fmt.Sprintf("amqp://%s:%s@localhost:5672/", userName, password)
	conn, err := amqp.Dial(connString)
	failOnError(err, "Failed to connect to RabbitMQ")
	rabbitMq = conn
}

func GetRabbitMq() *amqp.Connection {
	return rabbitMq
}
