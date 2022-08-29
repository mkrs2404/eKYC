package middlewares

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mkrs2404/eKYC/app/database"
	"github.com/mkrs2404/eKYC/app/models"
	"github.com/mkrs2404/eKYC/app/services"
	"gorm.io/gorm/logger"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error fetching the environment values")
	}
	database.Connect(os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_PORT"), logger.Silent)
	//Migrating tables to the database
	database.DB.AutoMigrate(&models.Plan{}, &models.Client{}, &models.File{}, &models.Api_Calls{})
	services.SeedPlanData()
	database.DB.Exec("DELETE FROM files")
	database.DB.Exec("DELETE FROM clients")

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestAuthMiddleware(t *testing.T) {

	gin.SetMode(gin.TestMode)
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)

	token, client, err := services.SetupClient()
	if err != nil {
		t.Fatal("Client setup failed")
	}
	ctx.Request, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(""))
	ctx.Request.Header.Set("Authorization", token)

	handlerFunc := AuthRequired()
	handlerFunc(ctx)

	clientInterface, _ := ctx.Get("client")
	authenticatedClient := clientInterface.(models.Client)

	if client.ID != authenticatedClient.ID || client.Name != authenticatedClient.Name || client.Email != authenticatedClient.Email || client.Plan != authenticatedClient.Plan {
		t.Errorf("Expected %v, Got %v", client, authenticatedClient)
	}

	database.DB.Exec("DELETE FROM clients")
}
