package controllers

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/mkrs2404/eKYC/api/resources"
	"github.com/mkrs2404/eKYC/helper"
	"github.com/mkrs2404/eKYC/messages"
	"github.com/mkrs2404/eKYC/minio_client"
)

const bucketName = "images"

func UploadImageClient(c *gin.Context) {

	clientEmail := c.GetString("client")
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

	fileName, err := saveFile(uploadImageRequest.Image)
	fmt.Println("filename ", fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.FILE_UPLOAD_FAILED,
		})
		c.Abort()
		return
	}

	ctx := context.Background()

	bucketExists, err := minio_client.Minio.BucketExists(ctx, bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.FILE_UPLOAD_FAILED,
		})
		c.Abort()
		return
	}
	if !bucketExists {
		err = minio_client.Minio.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errorMsg": messages.FILE_UPLOAD_FAILED,
			})
			c.Abort()
			return
		}
	}
	s3FileName := fmt.Sprintf("%s/%s", clientEmail, fileName)
	filePath := fmt.Sprintf("./uploads/%s", fileName)
	_, err = minio_client.Minio.FPutObject(ctx, bucketName, s3FileName, filePath, minio.PutObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.FILE_UPLOAD_FAILED,
		})
		fmt.Println(err)
		c.Abort()
		return
	}
	deleteLocalFile(filePath)

	c.JSON(http.StatusOK, gin.H{
		"msg":         "success",
		"clientEmail": clientEmail,
	})
}

func saveFile(image *multipart.FileHeader) (string, error) {

	var err error
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		return "", err
	}
	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().Unix(), filepath.Ext(image.Filename)))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	imageFile, _ := image.Open()
	_, err = io.Copy(dst, imageFile)
	return filepath.Base(dst.Name()), err
}

func deleteLocalFile(filePath string) {
	os.Remove(filePath)
}
