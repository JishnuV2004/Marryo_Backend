package models

import "gorm.io/gorm"

type Match struct {
	gorm.Model

	User1ID uint `gorm:"index;not null"`
	User2ID uint `gorm:"index;not null"`

	User1 User `gorm:"foreignKey:User1ID;references:ID"`
	User2 User `gorm:"foreignKey:User2ID;references:ID"`
}
