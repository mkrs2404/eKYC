package seeds

import "gorm.io/gorm"

type Seed struct {
	db *gorm.DB
}

func NewSeed(db *gorm.DB) Seed {
	return Seed{
		db: db,
	}
}
