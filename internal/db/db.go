package db

import (
	"fmt"

	"github.com/Tolfx/audioglyph/internal/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(dsn string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(
		&models.Artist{},
		&models.Release{},
		&models.PhysicalMedia{},
		&models.Track{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	return &Database{DB: db}, nil
}
