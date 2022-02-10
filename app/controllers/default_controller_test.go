package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type testData struct {
	testURL      string
	expectedCode int
	expectedBody string
}

//Welcome Page Data
var welcomeTestData = []testData{
	{
		testURL:      "/",
		expectedCode: 200,
		expectedBody: "Welcome to the eKYC API Portal.",
	},
}

//Invalid route data
var noRouteTestdata = []testData{
	{
		testURL:      "/random",
		expectedCode: 404,
		expectedBody: "Endpoint doesn't exist",
	},
}

func TestWelcomeController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/", WelcomePage)
	for _, data := range welcomeTestData {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest(http.MethodGet, data.testURL, nil)
		WelcomePage(ctx)
		if resRecorder.Code != data.expectedCode {
			t.Errorf("Expected %d, Got %d", data.expectedCode, resRecorder.Code)
		}
	}
}

func TestNoRouteController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.NoRoute(NoRoute)
	for _, data := range noRouteTestdata {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest(http.MethodGet, data.testURL, nil)
		NoRoute(ctx)
		if resRecorder.Code != data.expectedCode {
			t.Errorf("Expected %d, Got %d", data.expectedCode, resRecorder.Code)
		}
	}
}
