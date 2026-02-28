package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username   string `gorm:"default:user"`
	Email      string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	Role       string `gorm:"default:user"`
	Status     bool   `gorm:"default:false"`
	IsVerified bool   `gorm:"default:false"`
	Roles []Role `gorm:"many2many:user_roles;"`

	Profile Profile `gorm:"constraint:OnDelete:CASCADE;"`
}

