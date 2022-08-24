package controllers

import (
	"github.com/gin-gonic/gin"
)

// SignUpClient godoc
// @Summary  Signs up a client
// @ID       sign-up-client
// @Accept   json
// @Produce  json
// @Param    message  body      resources.SignUpRequest  true  "Client Info"
// @Success  200      {object}  object{access_key=string}
// @Failure  400      "Invalid Request"
// @Failure  500      "Internal Server Error"
// @Router   /signup [post]
//Handler for /api/v1/signup
func SignUpClient(c *gin.Context) {

	//Validate request

	//Save the client to the DB

	//Generate JWT token to send back as response

	//Enable these after implementing the above
	// c.Set("access_key", token)
	// c.JSON(http.StatusCreated, gin.H{
	// 	"access_key": token,
	// })

}
