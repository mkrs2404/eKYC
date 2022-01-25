package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/api/resources"
	"github.com/mkrs2404/eKYC/auth"
	"github.com/mkrs2404/eKYC/helper"
	"github.com/mkrs2404/eKYC/messages"
)

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

		clientEmail, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errorMsg": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("client", clientEmail)
		c.Next()
	}
}

func extractToken(header resources.AuthHeader) (string, error) {
	jwtToken := strings.Split(header.JWTToken, "Bearer ")
	var err error
	if len(jwtToken) < 2 {
		err = errors.New(messages.PROVIDE_PROPER_AUTH_HEADER)
		return "", err
	}
	return jwtToken[1], err
}
