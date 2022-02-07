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

var ocrTestData = []struct {
	body         string
	expectedCode int
	imgType      string
}{
	//Positive case
	{
		expectedCode: 200,
		imgType:      "id_card",
	},
	//UUID of face type
	{
		expectedCode: 400,
		imgType:      "face",
	},
	//Invalid UUID
	{
		body:         `{"image":"abcd"}`,
		expectedCode: 400,
		imgType:      "id_card",
	},
	//Missing UUID
	{
		body:         `{"image":""}`,
		expectedCode: 400,
		imgType:      "id_card",
	},
}

const ocrUrl = "/api/v1/ocr"

func TestOcrClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST(ocrUrl, middlewares.AuthRequired(), OcrClient)

	token, client, err := services.SetupClient()
	if err != nil {
		t.Fatal("Error setting up the client")
	}

	fileUUID1, fileUUID2, fileInfo1, fileInfo2 := services.SetupImageUpload(client, "id_card", "face")

	for _, data := range ocrTestData {

		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)

		ctx.Set("client", client)

		//Setting separate bucket name for tests
		ctx.Set("testBucket", "test")

		if data.body == "" {
			if data.imgType == "id_card" {
				data.body = fmt.Sprintf(`{"image":"%s"}`, fileUUID1)
			} else {
				data.body = fmt.Sprintf(`{"image":"%s"}`, fileUUID2)
			}
		}
		ctx.Request, _ = http.NewRequest(http.MethodPost, ocrUrl, strings.NewReader(data.body))

		//Setting the authorization token
		ctx.Request.Header.Set("Authorization", token)
		OcrClient(ctx)

		if resRecorder.Code != data.expectedCode {
			t.Errorf("Expected %d %s, Got %d %s", data.expectedCode, ctx.Request.Body, resRecorder.Code, resRecorder.Body)
		}
	}

	Clear(fileInfo1.Key, fileInfo2.Key)

}
