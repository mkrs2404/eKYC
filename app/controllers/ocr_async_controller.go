package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/app/helper"
	"github.com/mkrs2404/eKYC/app/messages"
	"github.com/mkrs2404/eKYC/app/models"
	"github.com/mkrs2404/eKYC/app/resources"
	"github.com/mkrs2404/eKYC/app/services"
	"syreclabs.com/go/faker"
)

func AsyncOcrClient(c *gin.Context) {

	const apiType = "ocr"

	//Getting the client object from previous http.handler
	clientInterface, _ := c.Get("client")
	client := clientInterface.(models.Client)

	//Binding the request to the model
	var ocrRequest resources.OcrRequest
	err := c.ShouldBindJSON(&ocrRequest)
	failure := helper.ReportValidationFailure(err, c)
	if failure {
		return
	}

	//Fetching the file details if it exists for the client
	file, err := services.GetFileForClient(ocrRequest.Image, client.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorMsg": messages.INVALID_IMAGE_ID,
		})
		c.Abort()
		return
	}

	//If image is not an id card
	if file.FileType != "id_card" {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorMsg": messages.NOT_AN_ID_CARD,
		})
		c.Abort()
		return
	}

	//Saving the api call info into the DB
	request, _ := json.Marshal(ocrRequest)
	apiCall, err := services.SaveApiCall(request, nil, -1, apiType, client.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.DATABASE_SAVE_FAILED,
		})
		c.Abort()
		return
	}

	//goroutine mimicking ML workload
	go ocrWorker(apiCall)

	c.JSON(http.StatusOK, gin.H{
		"job_id": apiCall.ID,
	})

}

func GetOcrData(c *gin.Context) {

	//Getting the client object from previous http.handler
	clientInterface, _ := c.Get("client")
	client := clientInterface.(models.Client)

	//Binding the request to the model
	var ocrJob resources.JobRequest
	err := c.ShouldBindJSON(&ocrJob)
	failure := helper.ReportValidationFailure(err, c)
	if failure {
		return
	}
	err = services.ValidateJobId(ocrJob.JobId, int(client.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorMsg": messages.INVALID_IMAGE_ID,
		})
		c.Abort()
		return
	}

	key := strconv.Itoa(ocrJob.JobId)
	jsonPayload, err := services.GetFromRedis(key)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "still processing",
		})
		c.Abort()
		return
	}

	ocrData := resources.OcrData{}
	json.Unmarshal([]byte(jsonPayload), &ocrData)

	c.JSON(http.StatusOK, ocrData)
}

func ocrWorker(apiCall models.Api_Calls) {

	//Simulating ML workload
	time.Sleep(10 * time.Second)

	//Adding fake data
	ocrData := resources.OcrData{
		Name:         faker.Name().Name(),
		Gender:       "male",
		DateOfBirth:  faker.Date().Birthday(18, 80).String()[:10],
		IdNumber:     strconv.FormatInt(faker.Number().NumberInt64(12), 10),
		AddressLine1: fmt.Sprintf("%s, %s,", faker.Address().SecondaryAddress(), faker.Address().StreetAddress()),
		AddressLine2: faker.Address().City(),
		Pincode:      faker.Address().Postcode(),
	}

	json, err := json.Marshal(ocrData)
	if err != nil {
		log.Fatal(err)
	}
	//Setting the score in Redis
	err = services.SetToRedis(strconv.Itoa(int(apiCall.ID)), json, time.Hour)
	if err != nil {
		log.Fatal(err)
	}
}
