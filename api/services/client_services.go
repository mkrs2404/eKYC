package services

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/mkrs2404/eKYC/api/models"
	"github.com/mkrs2404/eKYC/api/resources"
	"github.com/mkrs2404/eKYC/database"
)

/*SaveClient saves the client to the DB and returns back the client and an error(if there's email duplicacy)*/
func SaveClient(signUpRequest resources.SignUpRequest) (models.Client, error) {

	//Create a client object from the provided request data
	var client models.Client
	client.Name = signUpRequest.Name
	client.Email = signUpRequest.Email
	client.Plan = GetPlanId(signUpRequest.Plan)

	err := database.DB.Create(&client).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); pgErr.Code == "23505" {
			err = errors.New("client already exists")
		}
	}
	return client, err
}
