package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/mkrs2404/eKYC/api/middlewares"
	"github.com/mkrs2404/eKYC/api/models"
	"github.com/mkrs2404/eKYC/api/services"
	"github.com/mkrs2404/eKYC/database"
	"github.com/mkrs2404/eKYC/minio_client"
)

var middlewareTestData = []struct {
	token        string
	expectedCode int
}{
	//Empty Authorization header
	{
		token:        "",
		expectedCode: http.StatusBadRequest,
	},
	//Incorrect Token
	{
		token:        "k",
		expectedCode: http.StatusUnauthorized,
	},
}
var imageUploadTestData = []struct {
	filePath     string
	imageType    string
	expectedCode int
}{
	//Positive case
	{
		filePath:     "../../test_assets/architecture.png",
		imageType:    "face",
		expectedCode: 200,
	},
	//Positive case
	{
		filePath:     "../../test_assets/architecture.png",
		imageType:    "id_card",
		expectedCode: 200,
	},
	//Incorrect file extension
	{
		filePath:     "../../test_assets/architecture.gif",
		imageType:    "face",
		expectedCode: 400,
	},
	//Incorrect file size
	{
		filePath:     "../../test_assets/huge.jpeg",
		imageType:    "face",
		expectedCode: 400,
	},
	//Incorrect image type
	{
		filePath:     "../../test_assets/architecture.png",
		imageType:    "passport",
		expectedCode: 400,
	},
	//Missing image type
	{
		filePath:     "../../test_assets/architecture.png",
		imageType:    "",
		expectedCode: 400,
	},
}

const imageUploadUrl = "/api/v1/image"

func TestImageUploadClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST(imageUploadUrl, middlewares.AuthRequired(), UploadImageClient)

	for _, data := range imageUploadTestData {

		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)

		multiWriter, body, err := createMultipartPayload(data.filePath, data.imageType)
		if err != nil {
			t.Fatal("Error setting up the image file")
		}

		ctx.Request, _ = http.NewRequest(http.MethodPost, imageUploadUrl, body)
		ctx.Request.Header.Set("Content-Type", multiWriter.FormDataContentType())

		token, err := setupClient(ctx)
		if err != nil {
			t.Fatal("Error setting up the client")
		}
		//Setting the authorization token
		ctx.Request.Header.Set("Authorization", token)

		//Setting separate bucket name for tests
		ctx.Set("testBucket", "test")
		UploadImageClient(ctx)

		router.ServeHTTP(resRecorder, ctx.Request)

		if resRecorder.Code != data.expectedCode {
			t.Errorf("Expected %d, Got %d ", data.expectedCode, resRecorder.Code)
		}

		//Clearing up test DB
		database.DB.Exec("DELETE FROM files")
		database.DB.Exec("DELETE FROM clients")

		objectName := ctx.GetString("filePath")
		//Removing the bucket name from the filePath
		if objectName != "" {
			objectName = strings.Split(objectName, services.BucketName+"/")[1]
		}

		//Deleting the test images uploaded to minio
		minio_client.Minio.RemoveObject(context.Background(), services.BucketName, objectName, minio.RemoveObjectOptions{})
	}

	//Removing the uploaded files
	os.Remove("./uploads")
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST(imageUploadUrl, middlewares.AuthRequired(), UploadImageClient)

	for _, data := range middlewareTestData {

		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		body := ""
		ctx.Request, _ = http.NewRequest(http.MethodPost, imageUploadUrl, strings.NewReader(body))

		token, err := setupClient(ctx)
		if err != nil {
			t.Fatal("Error setting up the client")
		}

		if data.token == "" {
			ctx.Request.Header.Set("Authorization", "")
		} else {
			ctx.Request.Header.Set("Authorization", token+data.token)
		}

		router.ServeHTTP(resRecorder, ctx.Request)

		if resRecorder.Code != data.expectedCode {
			t.Errorf("Expected %d, Got %d ", data.expectedCode, resRecorder.Code)
		}
		database.DB.Exec("DELETE FROM clients")
	}

}

//createMultipartPayload takes the local filepath and imagetype and generates a Multipart paylaod
func createMultipartPayload(filePath, imageType string) (*multipart.Writer, *bytes.Buffer, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}
	body := new(bytes.Buffer)
	multiWriter := multipart.NewWriter(body)

	w, err := multiWriter.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, nil, err
	}
	w.Write(fileContents)

	multiWriter.WriteField("type", imageType)

	multiWriter.Close()

	return multiWriter, body, err
}

//setupClient creates a client in DB and returns the Auth header for Image upload tests
func setupClient(ctx *gin.Context) (string, error) {

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST(signUpUrl, SignUpClient)
	signUpRes := httptest.NewRecorder()
	signUpCtx, _ := gin.CreateTestContext(signUpRes)
	signUpCtx.Request, _ = http.NewRequest(http.MethodPost, signUpUrl, strings.NewReader(`{"name": "bob","email": "bob@one2.in","plan": "basic"}`))

	SignUpClient(signUpCtx)

	router.ServeHTTP(signUpRes, signUpCtx.Request)

	var client models.Client
	err := database.DB.Where("email = ?", "bob@one2.in").First(&client).Error
	ctx.Set("client", client)

	//Creating the auth header
	token := fmt.Sprintf("Bearer %s", signUpCtx.GetString("access_key"))

	return token, err
}