package services

import (
	"context"
	"log"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/mkrs2404/eKYC/api/models"
)

//SetupImageUpload uploads the files to minio and saves their data in DB, to setup the environment for tests
//If the imgType2 is empty then only 1 file will be uploaded and saved
func SetupImageUpload(client models.Client, imgType1 string, imgType2 string) (uuid.UUID, uuid.UUID, minio.UploadInfo, minio.UploadInfo) {

	var err1, err2 error
	var fileInfo1, fileInfo2 minio.UploadInfo
	var fileUUID1, fileUUID2 uuid.UUID
	filePath := []string{"../../test_assets/db5ed785-e54e-4ffb-9ed5-0653aea87217.png", "../../test_assets/f928e240-42da-490c-a4c2-14aac382f03b.png"}
	fileName := []string{filepath.Base(filePath[0]), filepath.Base(filePath[1])}

	//Uploading the files to minio
	fileInfo1, err1 = UploadToMinio(client.ID, fileName[0], imgType1, filePath[0], context.Background(), "")
	if imgType2 != "" {
		fileInfo2, err2 = UploadToMinio(client.ID, fileName[1], imgType2, filePath[1], context.Background(), "")
	}

	if err1 != nil || err2 != nil {
		log.Fatal("Upload to minio failed")
	}

	//Saving file's metadata to the database
	fileUUID1, err1 = SaveFile(fileInfo1.Bucket, fileInfo1.Key, fileInfo1.Size, imgType1, client.ID)
	if imgType2 != "" {
		fileUUID2, err2 = SaveFile(fileInfo2.Bucket, fileInfo2.Key, fileInfo2.Size, imgType2, client.ID)
	}

	if err1 != nil || err2 != nil {
		log.Fatal("DB save failed")
	}

	return fileUUID1, fileUUID2, fileInfo1, fileInfo2
}

func SetupAsyncFaceMatch(image1 string, image2 string, apiType string, client models.Client) (models.Api_Calls, error) {

	//Checking if both the images exist under the same client
	file1, err1 := GetFileForClient(image1, client.ID)
	file2, err2 := GetFileForClient(image2, client.ID)
	if err1 != nil || err2 != nil {
		return models.Api_Calls{}, err1
	}

	ctx := context.Background()

	//Downloading the files from minio
	_, err1 = DownloadFromMinio(ctx, file1.FileStoragePath, file1.FileName)
	_, err2 = DownloadFromMinio(ctx, file2.FileStoragePath, file2.FileName)

	if err1 != nil || err2 != nil {
		return models.Api_Calls{}, err1
	}

	//Saving the api call info into the DB
	apiCall, err := SaveApiCall(-1, apiType, client.ID)
	if err != nil {
		return models.Api_Calls{}, err
	}

	return apiCall, nil
}

func SetupAsyncOcr(image string, apiType string, client models.Client) (models.Api_Calls, error) {

	//Checking if the image exists under the same client
	file, err := GetFileForClient(image, client.ID)
	if err != nil {
		return models.Api_Calls{}, err
	}

	ctx := context.Background()

	//Downloading the file from minio
	_, err = DownloadFromMinio(ctx, file.FileStoragePath, file.FileName)

	if err != nil {
		return models.Api_Calls{}, err
	}

	//Saving the api call info into the DB
	apiCall, err := SaveApiCall(-1, apiType, client.ID)
	if err != nil {
		return models.Api_Calls{}, err
	}

	return apiCall, nil
}
