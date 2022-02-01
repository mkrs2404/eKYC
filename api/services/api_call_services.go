package services

import (
	"github.com/mkrs2404/eKYC/api/models"
	"github.com/mkrs2404/eKYC/database"
)

//SaveApiCall saves the api calls into the database
func SaveApiCall(matchResult int, apiType string, clientId uint) error {

	var api_call models.Api_Calls

	api_call.MatchResult = float32(matchResult)
	api_call.Type = apiType
	api_call.ClientID = clientId

	err := database.DB.Create(&api_call).Error
	return err
}
