package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkrs2404/eKYC/app/messages"
	"github.com/mkrs2404/eKYC/app/services"
)

func UpdateApi() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

		api_call_id, _ := c.Get("api_call_id")
		response, _ := c.Get("response")
		responseStatus := c.Writer.Status()

		if response == nil {
			return
		}
		_, err := services.UpdateApiCallResponse(api_call_id.(int), response.([]byte), responseStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errorMsg": messages.DATABASE_SAVE_FAILED,
			})
			c.Abort()
			return
		}
	}
}
