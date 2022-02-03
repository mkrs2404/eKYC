package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/api/models"
	"github.com/mkrs2404/eKYC/api/resources"
	"github.com/mkrs2404/eKYC/api/services"
	"github.com/mkrs2404/eKYC/helper"
	"github.com/mkrs2404/eKYC/messages"
	"syreclabs.com/go/faker"
)

func OcrClient(c *gin.Context) {

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

	//Adding fake data
	name := faker.Name().Name()
	gender := "male"
	dateOfBirth := faker.Date().Birthday(18, 80).String()[:10]
	idNumber := strconv.FormatInt(faker.Number().NumberInt64(12), 10)
	addressLine1 := fmt.Sprintf("%s, %s,", faker.Address().SecondaryAddress(), faker.Address().StreetAddress())
	addressLine2 := faker.Address().City()
	pincode := faker.Address().Postcode()

	//Saving the api call info into the DB
	err = services.SaveApiCall(-1, apiType, client.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.DATABASE_SAVE_FAILED,
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":         name,
		"gender":       gender,
		"dateOfBirth":  dateOfBirth,
		"idNumber":     idNumber,
		"addressLine1": addressLine1,
		"addressLine2": addressLine2,
		"pincode":      pincode,
	})

}
