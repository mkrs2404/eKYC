package services

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/mkrs2404/eKYC/minio_client"
)

var BucketName = "images"

func CreateBucket(ctx context.Context) error {

	//Checking if the bucket exists in minio. If not, then creating a bucket
	bucketExists, err := minio_client.Minio.BucketExists(ctx, BucketName)
	if err != nil {
		return err
	}
	if !bucketExists {
		err = minio_client.Minio.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{})
	}
	return err
}

func UploadToMinio(clientId uint, fileName string, imageType string, filePath string, ctx context.Context, testBucketName string) (minio.UploadInfo, error) {

	//Creating folder structure for s3 bucket as bucketName -> 12 -> face -> fileName
	s3FileName := fmt.Sprintf("%d/%s/%s", clientId, imageType, fileName)

	if testBucketName != "" {
		BucketName = testBucketName
	}

	//Storing the image in minio
	fileInfo, err := minio_client.Minio.FPutObject(ctx, BucketName, s3FileName, filePath, minio.PutObjectOptions{})
	return fileInfo, err
}
