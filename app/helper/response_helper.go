package helper

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

/*ReportValidationFailure sends the response back with proper errors during validation.
It returns whether there was a validation error or not*/
func ReportValidationFailure(err error, c *gin.Context) bool {
	if err == nil {
		return false
	}
	var validatorErr validator.ValidationErrors
	var errorMsg []string
	//Checking the type of validation error
	if errors.As(err, &validatorErr) {
		for _, error := range validatorErr {
			errorMsg = append(errorMsg, error.Field()+" : "+MsgForTag(error.Tag(), error.Field()))
		}
		//Sending the error response
		c.JSON(http.StatusBadRequest, gin.H{
			"errorMsg": errorMsg,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorMsg": err.Error(),
		})
	}
	c.Abort()
	return true
}

//MsgForTag returns the error message for the type of error passed to it
func MsgForTag(tag string, field string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid"
	case "oneof":
		if strings.EqualFold(field, "imagetype") {
			return "Should be one of face or id_card"
		} else if strings.EqualFold(field, "plan") {
			return "Should be one of basic, advanced or enterprise"
		}
		return ""
	}
	return ""
}
