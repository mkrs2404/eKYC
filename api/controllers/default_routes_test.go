package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

var defaultTestData = []struct {
	testURL      string
	expectedCode int
	expectedBody string
}{
	{
		testURL:      "/",
		expectedCode: 200,
		expectedBody: "Welcome to the eKYC API Portal.",
	},
	{
		testURL:      "/random",
		expectedCode: 404,
		expectedBody: "Endpoint doesn't exist",
	},
}

func TestDefaultRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	for _, data := range defaultTestData {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest(http.MethodPost, data.testURL, nil)
		if data.testURL == "/" {
			WelcomePage(ctx)
		} else {
			NoRoute(ctx)
		}
		router.ServeHTTP(resRecorder, ctx.Request)
		if resRecorder.Code != data.expectedCode {
			t.Errorf("Expected %d, Got %d", data.expectedCode, resRecorder.Code)
		}
	}
}
