package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email      string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	IsVerified bool   `gorm:"default:false"`

	Profile Profile `gorm:"constraint:OnDelete:CASCADE;"`
}
