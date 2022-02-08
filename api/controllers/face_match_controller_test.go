package controllers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
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

		//Setting separate bucket name for tests
		ctx.Set("testBucket", "test")

		fileUUID1, fileUUID2, fileInfo1, fileInfo2 := services.SetupImageUpload(client, "face", "face")

		if data.body == "" {
			data.body = fmt.Sprintf(`{"image1":"%s", "image2":"%s"}`, fileUUID1, fileUUID2)
		}
		ctx.Request, _ = http.NewRequest(http.MethodPost, faceMatchUrl, strings.NewReader(data.body))

		//Setting the authorization token
		ctx.Request.Header.Set("Authorization", token)
		FaceMatchClient(ctx)

		if resRecorder.Code != data.expectedCode {
			t.Errorf("Expected %d, Got %d ", data.expectedCode, resRecorder.Code)
		}

		Clear(fileInfo1.Key, fileInfo2.Key)
	}
}

func Clear(objectName1, objectName2 string) {
	//Clearing up test DB
	database.DB.Exec("DELETE FROM api_calls")
	database.DB.Exec("DELETE FROM files")
	database.DB.Exec("DELETE FROM clients")

	//Deleting the test images uploaded to minio
	minio_client.Minio.RemoveObject(context.Background(), services.BucketName, objectName1, minio.RemoveObjectOptions{})
	if objectName2 != "" {
		minio_client.Minio.RemoveObject(context.Background(), services.BucketName, objectName2, minio.RemoveObjectOptions{})
	}
}
