// internal/models/device_token.go
package models

import "gorm.io/gorm"

type DeviceToken struct {
	gorm.Model
	UserID   uint   `gorm:"index;not null"`
	Token    string `gorm:"unique;not null"`
	Platform string // android | ios | web
}
