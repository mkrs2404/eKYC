package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/app/auth"
	"github.com/mkrs2404/eKYC/app/helper"
	"github.com/mkrs2404/eKYC/app/messages"
	"github.com/mkrs2404/eKYC/app/resources"
	"github.com/mkrs2404/eKYC/app/services"
	"gorm.io/gorm"
)

//AuthRequired is the middleware to authenticate the JWT token supplied in the header
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var header resources.AuthHeader

		err := c.ShouldBindHeader(&header)
		failure := helper.ReportValidationFailure(err, c)
		if failure {
			return
		}

		tokenString, err := extractToken(header)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errorMsg": err.Error(),
			})
			c.Abort()
			return
		}

		//Validating if the JWT token provided is authentic
		clientId, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errorMsg": err.Error(),
			})
			c.Abort()
			return
		}

		//Setting the client object to the context for the next http.handler, when the token has been authenticated
		client, err := services.GetClient(clientId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{
					"errorMsg": messages.CLIENT_NOT_FOUND,
				})
				c.Abort()
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"errorMsg": messages.DB_FAILURE,
				})
				c.Abort()
				return
			}
		}
		c.Set("client", client)
		c.Next()
	}
}

//extractToken extracts the JWT token from the Authorization header
func extractToken(header resources.AuthHeader) (string, error) {
	jwtToken := strings.Split(header.JWTToken, "Bearer ")
	var err error
	if len(jwtToken) < 2 {
		err = errors.New(messages.PROVIDE_PROPER_AUTH_HEADER)
		return "", err
	}
	return jwtToken[1], err
}
