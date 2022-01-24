package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WelcomePage(c *gin.Context) {
	c.JSON(http.StatusOK, "Welcome to the eKYC API Portal")
}

func NoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, "Endpoint doesn't exist")
}
