package workers

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/mkrs2404/eKYC/app/rabbitmq"
	"github.com/mkrs2404/eKYC/app/resources"
	"github.com/mkrs2404/eKYC/app/services"
	"github.com/streadway/amqp"
)

func FaceMatchWorker(jobs <-chan amqp.Delivery) {

	for job := range jobs {

		var faceMatchJob resources.FaceMatchJob
		json.Unmarshal(job.Body, &faceMatchJob)

		log.Println("Received : ", faceMatchJob)
		//Simulating ML workload
		time.Sleep(10 * time.Second)

		rand.Seed(time.Now().UnixNano())
		//Random score generation between 0-100
		faceMatchScore := rand.Intn(101)

		//Setting the score in Redis
		err := services.SetToRedis(strconv.Itoa(int(faceMatchJob.JobId)), faceMatchScore, time.Hour)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func FaceMatchUtil() {
	forever := make(chan bool)
	jobs := rabbitmq.SetupFaceMatchConsumer()
	go FaceMatchWorker(jobs)
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
