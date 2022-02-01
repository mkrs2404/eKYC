package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/mkrs2404/eKYC/api/resources"
	"github.com/mkrs2404/eKYC/messages"
	"github.com/mkrs2404/eKYC/minio_client"
)

var BucketName = "images"

func CreateBucket(ctx context.Context, c *gin.Context) bool {

	//Checking if the bucket exists in minio. If not, then creating a bucket
	bucketExists, err := minio_client.Minio.BucketExists(ctx, BucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.FILE_UPLOAD_FAILED,
			"reason":   messages.BUCKET_CREATION_FAILED,
		})
		c.Abort()
		return false
	}
	if !bucketExists {
		err = minio_client.Minio.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errorMsg": messages.FILE_UPLOAD_FAILED,
				"reason":   messages.BUCKET_CREATION_FAILED,
			})
			c.Abort()
			return false
		}
	}

	return true
}

func UploadToMinio(clientId uint, fileName string, uploadImageRequest resources.UploadImageRequest, ctx context.Context, c *gin.Context) (minio.UploadInfo, string, error) {

	//Creating folder structure for s3 bucket as bucketName -> 12 -> face -> fileName
	s3FileName := fmt.Sprintf("%d/%s/%s", clientId, uploadImageRequest.ImageType, fileName)
	filePath := fmt.Sprintf("./uploads/%s", fileName)

	testBucketName := c.GetString("testBucket")
	if testBucketName != "" {
		BucketName = testBucketName
	}

	//Storing the image in minio
	fileInfo, err := minio_client.Minio.FPutObject(ctx, BucketName, s3FileName, filePath, minio.PutObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.FILE_UPLOAD_FAILED,
			"reason":   messages.MINIO_UPLOAD_FAILED,
		})
		c.Abort()
	}
	return fileInfo, filePath, err
}

func DownloadFromMinio(ctx context.Context, objectName string, localFilePath string) error {

	//Object name is like - "images/19/face/17e1de1a-6229-4ffc-8635-03d6ce28de6e.png"
	//Separating bucketName from objectName
	bucketName := strings.Split(objectName, "/")[0]
	objectName = strings.Split(objectName, bucketName+"/")[1]
	err := minio_client.Minio.FGetObject(ctx, bucketName, objectName, localFilePath, minio.GetObjectOptions{})
	return err
}
