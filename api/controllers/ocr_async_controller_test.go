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

var asyncOcrTestData = []struct {
	body         string
	expectedCode int
}{
	//Positive case
	{
		expectedCode: 200,
	},
	//Invalid UUID
	{
		body:         `{"image":"abcd"}`,
		expectedCode: 400,
	},
	//Missing UUID
	{
		body:         `{"image":""}`,
		expectedCode: 400,
	},
}

var asyncGetOcrTestData = []struct {
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

const asyncOcrUrl = "/api/v1/ocr-async"
const asyncGetOcrDataUrl = "/api/v1/get-ocr-data"

func TestAsyncOcrClient(t *testing.T) {

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST(asyncOcrUrl, middlewares.AuthRequired(), AsyncOcrClient)

	for _, data := range asyncOcrTestData {

		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)

		token, client, err := services.SetupClient()
		if err != nil {
			t.Fatal("Error setting up the client")
		}
		ctx.Set("client", client)

		//Setting separate bucket name for tests
		ctx.Set("testBucket", "test")

		fileUUID, _, fileInfo, _ := services.SetupImageUpload(client, "id_card", "")

		if data.body == "" {
			data.body = fmt.Sprintf(`{"image":"%s"}`, fileUUID)
		}

		ctx.Request, _ = http.NewRequest(http.MethodPost, asyncOcrUrl, strings.NewReader(data.body))

		//Setting the authorization token
		ctx.Request.Header.Set("Authorization", token)
		AsyncOcrClient(ctx)

		if resRecorder.Code != data.expectedCode {
			t.Errorf("Expected %d, Got %d ", data.expectedCode, resRecorder.Code)
		}

		Clear(fileInfo.Key, "")
	}
}

func TestAsyncGetOcrData(t *testing.T) {

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST(asyncGetOcrDataUrl, middlewares.AuthRequired(), GetOcrData)

	for _, data := range asyncGetOcrTestData {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)

		token, client, err := services.SetupClient()
		if err != nil {
			t.Fatal("Error setting up the client")
		}
		ctx.Set("client", client)

		//Setting separate bucket name for tests
		ctx.Set("testBucket", "test")

		fileUUID, _, fileInfo, _ := services.SetupImageUpload(client, "id_card", "")

		api_call, err := services.SetupAsyncOcr(fileUUID.String(), "ocr", client)
		if err != nil {
			t.Fatal("Error setting up the Async OCR")
		}
		if data.body == "" {
			data.body = fmt.Sprintf(`{"job_id":%d}`, api_call.ID)
		}
		ctx.Request, _ = http.NewRequest(http.MethodPost, asyncGetOcrDataUrl, strings.NewReader(data.body))

		//Setting the authorization token
		ctx.Request.Header.Set("Authorization", token)
		GetOcrData(ctx)

		if resRecorder.Code != data.expectedCode {
			t.Errorf("Expected %d, Got %d ", data.expectedCode, resRecorder.Code)
		}

		Clear(fileInfo.Key, "")
	}
}
