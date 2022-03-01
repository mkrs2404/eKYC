package services

import (
	"github.com/mkrs2404/eKYC/app/database"
	"github.com/mkrs2404/eKYC/app/models"
)

//SaveApiCall saves the api calls into the database
func SaveApiCall(request []byte, response []byte, responseStatus int, apiType string, clientId uint) (models.Api_Calls, error) {

	var api_call models.Api_Calls
	api_call.Request = request
	api_call.Type = apiType
	api_call.ClientID = clientId
	api_call.Response = response
	api_call.ReponseStatus = responseStatus
	err := database.DB.Create(&api_call).Error
	return api_call, err
}

func UpdateApiCallResponse(api_call_id int, response []byte, responseStatus int) (models.Api_Calls, error) {

	var api_call models.Api_Calls
	err := database.DB.First(&api_call, api_call_id).Error
	if err != nil {
		return api_call, err
	}
	api_call.ReponseStatus = responseStatus
	api_call.Response = response
	err = database.DB.Save(&api_call).Error
	return api_call, err
}

func ValidateJobId(jobId, clientId int) error {

	var api_call models.Api_Calls
	err := database.DB.Where("id = ? AND client_id = ?", jobId, clientId).First(&api_call).Error
	return err
}

func GetApiCall(jobId int) (models.Api_Calls, error) {
	var api_call models.Api_Calls
	err := database.DB.First(&api_call, jobId).Error
	return api_call, err
}
