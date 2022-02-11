package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/app/messages"
	"github.com/mkrs2404/eKYC/app/services"
)

func SaveApi() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

		apiType := c.GetString("apitype")
		clientId := c.GetUint("clientid")
		request, _ := c.Get("request")
		response, _ := c.Get("response")
		responseStatus := c.Writer.Status()

		if response == nil {
			return
		}
		_, err := services.SaveApiCall(request.([]byte), response.([]byte), responseStatus, apiType, clientId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errorMsg": messages.DATABASE_SAVE_FAILED,
			})
			c.Abort()
			return
		}
	}
}
