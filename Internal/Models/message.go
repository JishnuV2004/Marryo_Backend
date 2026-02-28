package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model

	MatchID uint `gorm:"index"`
	SenderID uint
	Content  string
}
