package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/api/resources"
)

func (server *Server) signUp(c *gin.Context) {

	var signUpRequest resources.SignUpRequest
	err := c.ShouldBindJSON(&signUpRequest)
	reportValidationFailure(err, c)

	err = signUpRequest.Validate()
	reportValidationFailure(err, c)

}

func reportValidationFailure(err error, c *gin.Context) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Validation Failure", "error": err.Error()})
		c.Abort()
		return
	}
}
