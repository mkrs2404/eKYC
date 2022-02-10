package models

import "gorm.io/gorm"

type Api_Calls struct {
	gorm.Model
	Type        string
	MatchResult float32
	ClientID    uint
}
