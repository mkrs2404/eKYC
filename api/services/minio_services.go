package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/mkrs2404/eKYC/api/resources"
	"github.com/mkrs2404/eKYC/messages"
	"github.com/mkrs2404/eKYC/minio_client"
)

const bucketName = "images"

func CreateBucket(ctx context.Context, c *gin.Context) bool {

	//Checking if the bucket exists in minio. If not, then creating a bucket
	bucketExists, err := minio_client.Minio.BucketExists(ctx, bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.FILE_UPLOAD_FAILED,
			"reason":   messages.BUCKET_CREATION_FAILED,
		})
		c.Abort()
		return false
	}
	if !bucketExists {
		err = minio_client.Minio.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
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

	//Storing the image in minio
	fileInfo, err := minio_client.Minio.FPutObject(ctx, bucketName, s3FileName, filePath, minio.PutObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMsg": messages.FILE_UPLOAD_FAILED,
			"reason":   messages.MINIO_UPLOAD_FAILED,
		})
		c.Abort()
	}
	return fileInfo, filePath, err
}
