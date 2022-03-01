package rabbitmq

import (
	"os"

	"github.com/streadway/amqp"
)

func SetupFaceMatchConsumer() <-chan amqp.Delivery {

	conn := GetRabbitMq()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	var queueName = "face_queue"
	if os.Getenv("FACE_WORKER_QUEUE") != "" {
		queueName = os.Getenv("FACE_WORKER_QUEUE")
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	failOnError(err, "Failed to declare a queue")

	jobs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consume")

	return jobs
}
