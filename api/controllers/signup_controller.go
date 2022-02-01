package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/api/models"
	"github.com/mkrs2404/eKYC/api/resources"
	"github.com/mkrs2404/eKYC/api/services"
	"github.com/mkrs2404/eKYC/auth"
	"github.com/mkrs2404/eKYC/database"
	"github.com/mkrs2404/eKYC/helper"
	"github.com/mkrs2404/eKYC/messages"
)

var signUpUrl = "/api/v1/signup"

//Handler for /api/v1/signup
func SignUpClient(c *gin.Context) {

	var signUpRequest resources.SignUpRequest

	//Validating request
	err := c.ShouldBindJSON(&signUpRequest)
	failure := helper.ReportValidationFailure(err, c)
	if failure {
		return
	}

	err = signUpRequest.Validate()
	failure = helper.ReportValidationFailure(err, c)
	if failure {
		return
	}

	//Saving the client to the DB
	client, err := services.SaveClient(signUpRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": messages.SIGN_UP_FAILED,
			"err": err.Error(),
		})
		c.Abort()
		return
	}

	//Generating JWT token to send back as response
	token, err := auth.GenerateToken(client.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": messages.SIGN_UP_FAILED,
		})
		database.DB.Delete(models.Client{}, client.ID)
		c.Abort()
		return
	}

	c.Set("access_key", token)
	c.JSON(http.StatusCreated, gin.H{
		"access_key": token,
	})

}

//SetupClient creates a client in DB and returns the Auth header for Image upload tests
func SetupClient(ctx *gin.Context) (string, error) {

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
