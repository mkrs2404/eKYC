package services

import (
	"github.com/mkrs2404/eKYC/api/database"
	"github.com/mkrs2404/eKYC/api/models"
)

//This method fetches the plan id using the plan name
func GetPlanId(planName string) uint {
	var plan models.Plan
	database.DB.Debug().Where("plan_name = ?", planName).First(&plan)
	return plan.ID
}
