package controllers

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/app/helper"
	"github.com/mkrs2404/eKYC/app/messages"
	"github.com/mkrs2404/eKYC/app/models"
	"github.com/mkrs2404/eKYC/app/resources"
	"github.com/mkrs2404/eKYC/app/services"
)

// FaceMatchClient godoc
// @Summary  Matches 2 face type images
// @ID       face-match-client
// @Accept   json
// @Produce  json
// @Param    Authorization  header    string   true  "Authentication header"
// @Param                             message  body  resources.FaceMatchRequest    true  "Match Request Info"
// @Success  200            {object}           object{score=int}
// @Failure  400            "Invalid Request"
// @Failure  500            "Internal Server Error"
// @Router   /face-match [post]
func FaceMatchClient(c *gin.Context) {

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

	rand.Seed(time.Now().UnixNano())
	//Random score generation between 0-100
	faceMatchScore := rand.Intn(101)

	//Saving the api call info into the DB
	request, _ := json.Marshal(faceMatchRequest)
	c.Set("apitype", apiType)
	c.Set("clientid", client.ID)
	c.Set("request", request)
	// _, err = services.SaveApiCall(request, apiType, client.ID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"errorMsg": messages.DATABASE_SAVE_FAILED,
	// 	})
	// 	c.Abort()
	// 	return
	// }
	responseBody := gin.H{
		"score": faceMatchScore,
	}
	response, _ := json.Marshal(responseBody)
	c.Set("response", response)
	c.JSON(http.StatusOK, responseBody)
}
