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
	"github.com/mkrs2404/eKYC/minio_client"
	"gorm.io/gorm/logger"
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

//Setting up DB connection, data seeding and Minio connection
func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error fetching the environment values")
	}
	database.Connect(os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_PORT"), logger.Silent)
	services.SeedPlanData()
	database.DB.Exec("DELETE FROM files")
	database.DB.Exec("DELETE FROM clients")

	minio_client.InitializeMinio(os.Getenv("TEST_MINIO_SERVER"), os.Getenv("TEST_MINIO_USER"), os.Getenv("TEST_MINIO_PWD"))
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestSignUpClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/api/v1/signup", SignUpClient)

	for _, data := range signupTestData {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest(http.MethodPost, signUpUrl, strings.NewReader(data.body))
		SignUpClient(ctx)
		router.ServeHTTP(resRecorder, ctx.Request)

		if resRecorder.Code != data.expectedCode {
			t.Errorf("Expected %d, Got %d", data.expectedCode, resRecorder.Code)
		}
	}

}
