package controllers

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mkrs2404/eKYC/api/models"
	"github.com/mkrs2404/eKYC/api/resources"
	"github.com/mkrs2404/eKYC/api/services"
	"github.com/mkrs2404/eKYC/helper"
	"github.com/mkrs2404/eKYC/messages"
)

//Handler for /api/v1/image
func UploadImageClient(c *gin.Context) {

	//Getting the client object from previous http.handler
	clientInterface, _ := c.Get("client")
	client := clientInterface.(models.Client)

	//Binding the incoming multipart req to the model
	var uploadImageRequest resources.UploadImageRequest
	err := c.ShouldBind(&uploadImageRequest)
	failure := helper.ReportValidationFailure(err, c)
	if failure {
		return
	}

	err = uploadImageRequest.Validate()
	failure = helper.ReportValidationFailure(err, c)
	if failure {
		return
	}

	//Validating file size and extension
	err = services.ValidateFile(uploadImageRequest.Image)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorMsg": messages.FILE_UPLOAD_FAILED,
			"reason":   err.Error(),
		})
		c.Abort()
		return
	}

	//Saving the file locally
	fileName, err := saveFileToDisk(uploadImageRequest.Image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.FILE_UPLOAD_FAILED,
			"reason":   messages.DISK_SAVE_FAILED,
		})
		c.Abort()
		return
	}

	ctx := context.Background()
	testBucketName := c.GetString("testBucket")

	//Creating a S3 bucket
	err = services.CreateBucket(ctx, testBucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.FILE_UPLOAD_FAILED,
			"reason":   messages.BUCKET_CREATION_FAILED,
		})
		c.Abort()
	}

	//path of the file to be uploaded
	filePath := fmt.Sprintf("./uploads/%s", fileName)

	//Uploading the file to minio
	fileInfo, err := services.UploadToMinio(client.ID, fileName, uploadImageRequest.ImageType, filePath, ctx, testBucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.FILE_UPLOAD_FAILED,
			"reason":   messages.MINIO_UPLOAD_FAILED,
		})
		c.Abort()
	}

	//Setting the objectName to delete later, if required
	c.Set("filePath", fileInfo.Key)

	//Saving file's metadata to the database
	fileUUID, err := services.SaveFile(fileInfo.Bucket, fileInfo.Key, fileInfo.Size, uploadImageRequest.ImageType, client.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.FILE_UPLOAD_FAILED,
			"reason":   messages.DATABASE_SAVE_FAILED,
		})
		c.Abort()
		return
	}

	//Deleting the locally saved file
	services.DeleteLocalFile(filePath)

	c.JSON(http.StatusOK, gin.H{
		"id": fileUUID,
	})
}

//saveFileToDisk saves the uploaded file to the local file system, and returns the saved file's paths
func saveFileToDisk(image *multipart.FileHeader) (string, error) {

	var err error
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		return "", err
	}
	//Creating a file with name derived from UUID and the file extension
	fileUUID := uuid.New()
	dst, err := os.Create(fmt.Sprintf("./uploads/%s%s", fileUUID, filepath.Ext(image.Filename)))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	imageFile, _ := image.Open()
	_, err = io.Copy(dst, imageFile)
	return filepath.Base(dst.Name()), err
}
