package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID              uuid.UUID `gorm:"primarykey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	FileName        string `gorm:"uniqueIndex;not null"`
	FileType        string
	FileExtension   string
	FileSizeKB      float64
	FileStoragePath string
	ClientID        uint
}
