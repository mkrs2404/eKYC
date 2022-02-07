package controllers

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/api/models"
	"github.com/mkrs2404/eKYC/api/resources"
	"github.com/mkrs2404/eKYC/api/services"
	"github.com/mkrs2404/eKYC/helper"
	"github.com/mkrs2404/eKYC/messages"
)

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
	apiCall, err := services.SaveApiCall(-1, apiType, client.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.DATABASE_SAVE_FAILED,
		})
		c.Abort()
		return
	}

	//goroutine mimicking ML workload
	go faceMatchWorker(apiCall)

	c.JSON(http.StatusOK, gin.H{
		"job_id": apiCall.ID,
	})

}

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

	c.JSON(http.StatusOK, gin.H{
		"score": score,
	})

}

func faceMatchWorker(apiCall models.Api_Calls) {

	//Simulating ML workload
	time.Sleep(10 * time.Second)

	rand.Seed(time.Now().UnixNano())
	//Random score generation between 0-100
	faceMatchScore := rand.Intn(101)

	//Saving the api call info into the DB
	_, err := services.UpdateApiCall(apiCall, faceMatchScore)
	if err != nil {
		log.Fatal(err)
	}

	//Setting the score in Redis
	err = services.SetToRedis(strconv.Itoa(int(apiCall.ID)), faceMatchScore, time.Hour)
	if err != nil {
		log.Fatal(err)
	}
}
