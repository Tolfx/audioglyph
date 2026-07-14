package models

import "gorm.io/gorm"

// TODO: Should PhysicalMedia only be a reference to the actual media and the content it resides?

// PhysicalMedia is a struct that represents a physical media item in the database. Such as a CD or Vinyl
type PhysicalMedia struct {
	gorm.Model
	ReleaseID     uint         `gorm:"not null"`
	Release       Release      `gorm:"foreignKey:ReleaseID"`
	Format        PhysicalType `gorm:"type:varchar(10);not null"`
	Barcode       string       `gorm:"type:varchar(255);not null"`
	CatalogNumber string       `gorm:"type:varchar(255)"`
	IssuedYear    int          `gorm:"type:int;not null"`
	InCollection  bool         `gorm:"type:bool;not null"`
	Tracks        []Track      `gorm:"foreignKey:PhysicalMediaID"`
}

type PhysicalType string

const (
	PhysicalTypeCD    PhysicalType = "CD"
	PhysicalTypeTape  PhysicalType = "Tape"
	PhysicalTypeVinyl PhysicalType = "Vinyl"
)
