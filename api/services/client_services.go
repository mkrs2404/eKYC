package services

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/mkrs2404/eKYC/api/models"
	"github.com/mkrs2404/eKYC/api/resources"
	"github.com/mkrs2404/eKYC/database"
)

//CreateClientFromRequest creates a client object from the provided request data
func CreateClientFromRequest(signUpRequest resources.SignUpRequest) models.Client {
	var client models.Client
	client.Name = signUpRequest.Name
	client.Email = signUpRequest.Email
	client.Plan = GetPlanId(signUpRequest.Plan)
	return client
}

/*SaveClient saves the client to the DB
  and returns back an error if there's email duplicacy*/
func SaveClient(client *models.Client) error {
	err := database.DB.Debug().Create(client).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); pgErr.Code == "23505" {
			err = errors.New("client already exists")
		}
	}
	return err
}
