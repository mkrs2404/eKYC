package services

import (
	"github.com/mkrs2404/eKYC/api/models"
	"github.com/mkrs2404/eKYC/database"
)

//SaveApiCall saves the api calls into the database
func SaveApiCall(matchResult int, apiType string, clientId uint) (models.Api_Calls, error) {

	var api_call models.Api_Calls
	api_call.MatchResult = float32(matchResult)
	api_call.Type = apiType
	api_call.ClientID = clientId
	err := database.DB.Create(&api_call).Error
	return api_call, err
}

func UpdateApiCall(api_call models.Api_Calls, matchResult int) (models.Api_Calls, error) {

	database.DB.First(&api_call)
	api_call.MatchResult = float32(matchResult)
	err := database.DB.Save(&api_call).Error
	return api_call, err
}

func ValidateMatchId(matchId, clientId int) error {

	var api_call models.Api_Calls
	err := database.DB.Where("id = ? AND client_id = ?", matchId, clientId).First(&api_call).Error
	return err
}
