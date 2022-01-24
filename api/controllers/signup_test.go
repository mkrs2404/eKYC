package controllers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mkrs2404/eKYC/api/services"
	"github.com/mkrs2404/eKYC/database"
)

var signupTestData = []struct {
	expectedCode int
	body         string
}{
	//Valid request
	{
		body:         `{"name": "bob","email": "bob@one2.in","plan": "basic"}`,
		expectedCode: 201,
	},
	//Duplicate request
	{
		body:         `{"name": "bob","email": "bob@one2.in","plan": "basic"}`,
		expectedCode: 400,
	},
	//Invalid email
	{
		body:         `{"name": "bob","email": "bobone2.in","plan": "basic"}`,
		expectedCode: 400,
	},
	//Invalid plan
	{
		body:         `{"name": "bob","email": "bob@one2.in","plan": "secure"}`,
		expectedCode: 400,
	},
	//Missing plan
	{
		body:         `{"name": "bob","email": "bob@one2.in","plan": ""}`,
		expectedCode: 400,
	},
}

var signUpUrl = "/api/v1/signup"

func TestMain(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error fetching the environment values")
	} else {
		database.Connect(os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_PORT"))
		services.SeedPlanData()
		database.DB.Exec("DELETE FROM clients")
	}
}

func TestSignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	for _, data := range signupTestData {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest(http.MethodPost, signUpUrl, strings.NewReader(data.body))
		SignUp(ctx)
		router.ServeHTTP(resRecorder, ctx.Request)

		if resRecorder.Code != data.expectedCode {
			t.Errorf("Expected %d, Got %d %s", data.expectedCode, resRecorder.Code, resRecorder.Body.String())
		}

	}

}
