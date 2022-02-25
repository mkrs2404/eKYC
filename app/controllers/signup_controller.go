package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/app/auth"
	"github.com/mkrs2404/eKYC/app/database"
	"github.com/mkrs2404/eKYC/app/helper"
	"github.com/mkrs2404/eKYC/app/messages"
	"github.com/mkrs2404/eKYC/app/models"
	"github.com/mkrs2404/eKYC/app/resources"
	"github.com/mkrs2404/eKYC/app/services"
)

// SignUpClient godoc
// @Summary  Signs up a client
// @ID       sign-up-client
// @Accept   json
// @Produce  json
// @Param                   message  body  resources.SignUpRequest    true  "Client Info"
// @Success  200  {object}           object{access_key=string}
// @Failure  400  "Invalid Request"
// @Failure  500  "Internal Server Error"
// @Router   /signup [post]
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
