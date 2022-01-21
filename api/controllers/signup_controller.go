package controllers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mkrs2404/eKYC/api/resources"
	"github.com/mkrs2404/eKYC/api/services"
	"github.com/mkrs2404/eKYC/auth"
)

//Handler for /api/v1/signup
func SignUp(c *gin.Context) {

	//Fetching the JWT token delay from env variables
	tokenExpiryDelay := os.Getenv("TOKEN_EXPIRY_DELAY")
	if len(tokenExpiryDelay) == 0 {
		log.Fatal("Token expiry delay not found")
	}

	expiryDelay, err := strconv.Atoi(tokenExpiryDelay)
	if err != nil {
		log.Fatal("Incorrect delay provided")
	}

	var signUpRequest resources.SignUpRequest
	//Validating request
	err = c.ShouldBindJSON(&signUpRequest)
	failure := reportValidationFailure(err, c)
	if failure {
		return
	}

	err = signUpRequest.Validate()
	failure = reportValidationFailure(err, c)
	if failure {
		return
	}

	//Creating client object from the request
	client := services.CreateClientFromRequest(signUpRequest)

	//Saving the client to the DB
	err = services.SaveClient(&client)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Sign Up failed", "err": err.Error()})
		c.Abort()
		return
	}

	//Generating JWT token to send back as response
	token, err := auth.GenerateToken(int(client.ID), expiryDelay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Sign Up failed"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_key": token})
}

/*reportValidationFailure sends the response back with proper errors during validation.
  It returns whether there was a validation error or not*/
func reportValidationFailure(err error, c *gin.Context) bool {
	if err != nil {
		var validatorErr validator.ValidationErrors
		var errorMsg string
		//Checking the type of validation error
		if errors.As(err, &validatorErr) {
			for _, error := range validatorErr {
				errorMsg += error.Field() + " : " + msgForTag(error.Tag()) + "; "
			}
			//Sending the error response
			c.JSON(http.StatusBadRequest, gin.H{"errorMsg": errorMsg[:len(errorMsg)-2]})
		}
		c.Abort()
		return true
	}
	return false
}

//msgForTag returns the error message for the type of error passed to it
func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid"
	case "oneof":
		return "Should be one of basic advanced enterprise"
	}
	return ""
}
