package models

import "gorm.io/gorm"

//Model for api/v1/signup
type Client struct {
	gorm.Model
	Name      string `gorm:"type:varchar(100)" json:"name"`
	Email     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	AccessKey string `gorm:"type:varchar(100)" json:"access_key"`
	Plan      uint
	File      File `gorm:"foreignKey:ClientID"`
}
