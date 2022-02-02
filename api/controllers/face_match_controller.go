package controllers

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/api/models"
	"github.com/mkrs2404/eKYC/api/resources"
	"github.com/mkrs2404/eKYC/api/services"
	"github.com/mkrs2404/eKYC/helper"
	"github.com/mkrs2404/eKYC/messages"
)

const apiType = "face-match"

func FaceMatchClient(c *gin.Context) {

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
	filePath1, err1 := services.DownloadFromMinio(ctx, file1.FileStoragePath, file1.FileName)
	filePath2, err2 := services.DownloadFromMinio(ctx, file2.FileStoragePath, file2.FileName)

	if err1 != nil && err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.MINIO_DOWNLOAD_FAILED,
		})
		c.Abort()
		return
	}

	//Random score generation between 0-100
	faceMatchScore := rand.Intn(101)

	//Saving the api call info into the DB
	err = services.SaveApiCall(faceMatchScore, apiType, client.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.DATABASE_SAVE_FAILED,
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"score": faceMatchScore,
	})

	//Deleting the downloaded files from minio
	services.DeleteLocalFile(filePath1)
	services.DeleteLocalFile(filePath2)
}