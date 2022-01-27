package models

import (
	"time"

	"github.com/google/uuid"
)

//Model for /api/v1/image
type File struct {
	ID              uuid.UUID `gorm:"primary_key; unique; type:uuid;"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	FileName        string `gorm:"uniqueIndex;not null"`
	FileType        string
	FileExtension   string
	FileSizeKB      float64
	FileStoragePath string
	ClientID        uint
}
