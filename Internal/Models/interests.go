package models

import "gorm.io/gorm"

type Interest struct {
	gorm.Model

	SenderID   uint `gorm:"index;not null"`
	ReceiverID uint `gorm:"index;not null"`

	Status string `gorm:"type:varchar(20);default:'pending'"`
	// pending | accepted | rejected
	Sender   User `gorm:"foreignKey:SenderID;references:ID" json:"Sender"`
	Receiver User `gorm:"foreignKey:ReceiverID;references:ID" json:"Receiver"`

	// Sender Profile `gorm:"foreignKey:SenderID;references:UserID"`
}
