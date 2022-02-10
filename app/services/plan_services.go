package services

import (
	"github.com/mkrs2404/eKYC/app/database"
	"github.com/mkrs2404/eKYC/app/models"
)

//GetPlanId fetches the plan id using the plan name
func GetPlanId(planName string) uint {
	var plan models.Plan
	database.DB.Where("plan_name = ?", planName).First(&plan)
	return plan.ID
}
