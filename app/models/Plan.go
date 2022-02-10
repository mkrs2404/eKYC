package models

import "gorm.io/gorm"

//Model for api/v1/signup
type Plan struct {
	gorm.Model
	PlanName      string  `gorm:"type:varchar(20);uniqueIndex;not null" json:"plan_name"`
	DailyBaseCost float64 `json:"daily_base_cost"`
	ApiCost       float64 `json:"api_cost"`
	StorageCost   float64 `json:"storage_cost"`
	Client        Client  `gorm:"foreignKey:Plan"`
}
