package models

import "gorm.io/gorm"

type Release struct {
	gorm.Model
	ArtistID       uint            `gorm:"not null"`
	Artist         Artist          `gorm:"foreignKey:ArtistID"`
	Title          string          `gorm:"type:varchar(255);not null"`
	Year           int             `gorm:"type:int;not null"`
	PhysicalMedias []PhysicalMedia `gorm:"foreignKey:ReleaseID"`
}
