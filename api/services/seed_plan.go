package services

import (
	"github.com/mkrs2404/eKYC/api/models"
	"github.com/mkrs2404/eKYC/database"
)

//Plan modifications can be done here
var plans = []models.Plan{
	{
		PlanName:      "basic",
		DailyBaseCost: 10,
		ApiCost:       0.1,
		StorageCost:   0.1,
	},
	{
		PlanName:      "advanced",
		DailyBaseCost: 15,
		ApiCost:       0.05,
		StorageCost:   0.05,
	},
	{
		PlanName:      "enterprise",
		DailyBaseCost: 20,
		ApiCost:       0.1,
		StorageCost:   0.01,
	},
}

//SeedPlanData seeds the DB with the available plans.
func SeedPlanData() {
	localPlans := make([]models.Plan, len(plans))
	copy(localPlans, plans)

	for i, plan := range localPlans {
		database.DB.Where(models.Plan{PlanName: plan.PlanName}).FirstOrCreate(&plan)
		localPlans[i] = plan
	}

	for i, plan := range plans {
		database.DB.Debug().Model(&plan).Where("ID = ?", localPlans[i].ID).Updates(models.Plan{DailyBaseCost: plan.DailyBaseCost, ApiCost: plan.ApiCost, StorageCost: plan.StorageCost})
	}
}
