package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/api/middlewares"
	"github.com/mkrs2404/eKYC/api/services"
)

var asyncFaceMatchTestData = []struct {
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

var asyncGetScoreTestData = []struct {
	body         string
	expectedCode int
}{
	//Positive case
	{
		expectedCode: 200,
	},
	//Invalid job id
	{
		body:         `{"job_id":11}`,
		expectedCode: 400,
	},
	//Missing job id
	{
		body:         `{"job_id": }`,
		expectedCode: 400,
	},
}

const asyncFaceMatchUrl = "/api/v1/face-match-async"

const asyncGetFaceMatchScoreUrl = "/api/v1/get-score"

func TestAsyncFaceMatchClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST(asyncFaceMatchUrl, middlewares.AuthRequired(), AsyncFaceMatchClient)

	for _, data := range asyncFaceMatchTestData {
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
		ctx.Request, _ = http.NewRequest(http.MethodPost, asyncFaceMatchUrl, strings.NewReader(data.body))

		//Setting the authorization token
		ctx.Request.Header.Set("Authorization", token)
		AsyncFaceMatchClient(ctx)

		if resRecorder.Code != data.expectedCode {
			t.Errorf("Expected %d, Got %d ", data.expectedCode, resRecorder.Code)
		}

		Clear(fileInfo1.Key, fileInfo2.Key)
	}
}

func TestAsyncGetFaceMatchScore(t *testing.T) {

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST(asyncGetFaceMatchScoreUrl, middlewares.AuthRequired(), GetFaceMatchScore)

	for _, data := range asyncGetScoreTestData {
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

		api_call, err := services.SetupAsyncFaceMatch(fileUUID1.String(), fileUUID2.String(), "face-match", client)
		if err != nil {
			t.Fatal("Error setting up the Async face match")
		}
		if data.body == "" {
			data.body = fmt.Sprintf(`{"job_id":%d}`, api_call.ID)
		}
		ctx.Request, _ = http.NewRequest(http.MethodPost, asyncGetFaceMatchScoreUrl, strings.NewReader(data.body))

		//Setting the authorization token
		ctx.Request.Header.Set("Authorization", token)
		GetFaceMatchScore(ctx)

		if resRecorder.Code != data.expectedCode {
			t.Errorf("Expected %d, Got %d ", data.expectedCode, resRecorder.Code)
		}

		Clear(fileInfo1.Key, fileInfo2.Key)
	}
}
