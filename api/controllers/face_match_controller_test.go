package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/mkrs2404/eKYC/api/middlewares"
	"github.com/mkrs2404/eKYC/api/services"
	"github.com/mkrs2404/eKYC/database"
	"github.com/mkrs2404/eKYC/minio_client"
)

var faceMatchTestData = []struct {
	body         string
	expectedCode int
}{
	//Positive case
	{
		expectedCode: 200,
	},
	//Invalid UUID
	{
		body:         `{"image1":"abcd", "image2":"defg"}`,
		expectedCode: 400,
	},
	//Missing UUID
	{
		body:         `{"image1":"", "image2":"defg"}`,
		expectedCode: 400,
	},
}

const faceMatchUrl = "/api/v1/face-match"

func TestFaceMatchClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST(faceMatchUrl, middlewares.AuthRequired(), FaceMatchClient)

	for _, data := range faceMatchTestData {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)

		token, client, err := services.SetupClient()
		if err != nil {
			t.Fatal("Error setting up the client")
		}
		ctx.Set("client", client)

		filePath1 := "../../test_assets/db5ed785-e54e-4ffb-9ed5-0653aea87217.png"
		fileName1 := filepath.Base(filePath1)
		filePath2 := "../../test_assets/f928e240-42da-490c-a4c2-14aac382f03b.png"
		fileName2 := filepath.Base(filePath2)

		//Setting separate bucket name for tests
		ctx.Set("testBucket", "test")

		//Uploading the files to minio
		fileInfo1, err1 := services.UploadToMinio(client.ID, fileName1, "face", filePath1, context.Background(), "")
		fileInfo2, err2 := services.UploadToMinio(client.ID, fileName2, "face", filePath2, context.Background(), "")
		if err1 != nil || err2 != nil {
			log.Fatal("Upload to minio failed")
		}

		//Saving file's metadata to the database
		fileUUID1, err1 := services.SaveFile(fileInfo1.Bucket, fileInfo1.Key, fileInfo1.Size, "face", client.ID)
		fileUUID2, err2 := services.SaveFile(fileInfo2.Bucket, fileInfo2.Key, fileInfo2.Size, "face", client.ID)
		if err1 != nil || err2 != nil {
			log.Fatal("DB save failed")
		}

		if data.body == "" {
			data.body = fmt.Sprintf(`{"image1":"%s", "image2":"%s"}`, fileUUID1, fileUUID2)
		}
		ctx.Request, _ = http.NewRequest(http.MethodPost, faceMatchUrl, strings.NewReader(data.body))

		//Setting the authorization token
		ctx.Request.Header.Set("Authorization", token)
		FaceMatchClient(ctx)

		router.ServeHTTP(resRecorder, ctx.Request)
		if resRecorder.Code != data.expectedCode {
			t.Errorf("Expected %d, Got %d ", data.expectedCode, resRecorder.Code)
		}

		//Clearing up test DB
		database.DB.Exec("DELETE FROM api_calls")
		database.DB.Exec("DELETE FROM files")
		database.DB.Exec("DELETE FROM clients")

		//Deleting the test images uploaded to minio
		minio_client.Minio.RemoveObject(context.Background(), services.BucketName, fileInfo1.Key, minio.RemoveObjectOptions{})
		minio_client.Minio.RemoveObject(context.Background(), services.BucketName, fileInfo2.Key, minio.RemoveObjectOptions{})
	}
}
