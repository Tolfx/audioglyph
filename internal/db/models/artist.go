package models

import "gorm.io/gorm"

type Artist struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);not null"`
}
