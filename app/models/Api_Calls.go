package models

import (
	"time"

	"gorm.io/datatypes"
)

type Api_Calls struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	Type          string
	Request       datatypes.JSON `json:"request"`
	Response      datatypes.JSON `json:"response"`
	ReponseStatus int            `json:"response_status"`
	ClientID      uint
}
