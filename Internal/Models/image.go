package models

import "gorm.io/gorm"

type Img struct {
	gorm.Model

	ProfileID uint   `gorm:"index;not null"` // FK → profiles.id
	URL       string `gorm:"not null"`

	IsPrimary  bool `gorm:"default:false"`
	IsApproved bool `gorm:"default:false"`

	Order int
}