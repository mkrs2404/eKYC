package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/app/helper"
	"github.com/mkrs2404/eKYC/app/messages"
	"github.com/mkrs2404/eKYC/app/models"
	"github.com/mkrs2404/eKYC/app/rabbitmq"
	"github.com/mkrs2404/eKYC/app/resources"
	"github.com/mkrs2404/eKYC/app/services"
	"github.com/streadway/amqp"
)

// AsyncFaceMatchClient godoc
// @Summary  Matches 2 face type images(async)
// @ID       face-match-async-client
// @Accept   json
// @Produce  json
// @Param    Authorization  header    string                      true  "Authentication header"
// @Param    message        body      resources.FaceMatchRequest  true  "Match Request Info"
// @Success  200            {object}  object{job_id=int}
// @Failure  400            "Invalid Request"
// @Failure  500            "Internal Server Error"
// @Router   /face-match-async [post]
func AsyncFaceMatchClient(c *gin.Context) {

	const apiType = "face-match"

	//Getting the client object from previous http.handler
	clientInterface, _ := c.Get("client")
	client := clientInterface.(models.Client)

	//Binding the request to the model
	var faceMatchRequest resources.FaceMatchRequest
	err := c.ShouldBindJSON(&faceMatchRequest)
	failure := helper.ReportValidationFailure(err, c)
	if failure {
		return
	}

	//Checking if both the images exist under the same client
	file1, err1 := services.GetFileForClient(faceMatchRequest.Image1, client.ID)
	file2, err2 := services.GetFileForClient(faceMatchRequest.Image2, client.ID)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorMsg": messages.INVALID_IMAGE_ID,
		})
		c.Abort()
		return
	}

	ctx := context.Background()

	//Downloading the files from minio
	_, err1 = services.DownloadFromMinio(ctx, file1.FileStoragePath, file1.FileName)
	_, err2 = services.DownloadFromMinio(ctx, file2.FileStoragePath, file2.FileName)

	if err1 != nil && err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.MINIO_DOWNLOAD_FAILED,
		})
		c.Abort()
		return
	}

	//Saving the api call info into the DB
	request, _ := json.Marshal(faceMatchRequest)
	apiCall, err := services.SaveApiCall(request, nil, -1, apiType, client.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.DATABASE_SAVE_FAILED,
		})
		c.Abort()
		return
	}

	//Sending job to RabbitMQ
	ch, q := rabbitmq.SetupFaceMatchProducer()
	job := resources.CreateFaceMatchJob(faceMatchRequest.Image1, faceMatchRequest.Image2, apiCall.ID)
	SendJob(ch, q, job)

	c.JSON(http.StatusOK, gin.H{
		"job_id": apiCall.ID,
	})

}

func SendJob(ch *amqp.Channel, q amqp.Queue, job *resources.FaceMatchJob) {

	body, err := json.Marshal(job)
	failOnError(err, "Failed to publish a message")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	failOnError(err, "Failed to publish a message")
}

// GetFaceMatchScore godoc
// @Summary  Gets face match score
// @ID       get-face-match-score-client
// @Accept   json
// @Produce  json
// @Param    Authorization  header    string              true  "Authentication header"
// @Param    message        body      object{job_id=int}  true  "Job Info"
// @Success  200            {object}  object{score=int}
// @Failure  400            "Invalid Request"
// @Failure  500            "Internal Server Error"
// @Router   /get-score [post]
func GetFaceMatchScore(c *gin.Context) {

	//Getting the client object from previous http.handler
	clientInterface, _ := c.Get("client")
	client := clientInterface.(models.Client)

	//Binding the request to the model
	var faceMatchJob resources.JobRequest
	err := c.ShouldBindJSON(&faceMatchJob)
	failure := helper.ReportValidationFailure(err, c)
	if failure {
		return
	}
	err = services.ValidateJobId(faceMatchJob.JobId, int(client.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorMsg": messages.INVALID_IMAGE_ID,
		})
		c.Abort()
		return
	}

	key := strconv.Itoa(faceMatchJob.JobId)
	score, err := services.GetFromRedis(key)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "still processing",
		})
		c.Abort()
		return
	}

	responseBody := gin.H{
		"score": score,
	}
	response, _ := json.Marshal(responseBody)

	//Saving the api call info into the DB
	c.Set("api_call_id", faceMatchJob.JobId)
	c.Set("response", response)
	c.JSON(http.StatusOK, responseBody)

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
