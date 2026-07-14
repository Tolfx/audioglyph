package models

import "gorm.io/gorm"

type Track struct {
	gorm.Model
	PhysicalMediaID uint          `gorm:"not null"`
	PhysicalMedia   PhysicalMedia `gorm:"foreignKey:PhysicalMediaID"`
	Artist          string        `gorm:"type:varchar(255)"`
	Position        string        `gorm:"type:varchar(10)"`
	Title           string        `gorm:"type:varchar(255);not null"`
	DurationSeconds int           `gorm:"type:int;not null"`
}
