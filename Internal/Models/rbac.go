package models

import "gorm.io/gorm"

// stuct for role

type Role struct {
    gorm.Model
    Name        string       `gorm:"uniqueIndex"`
    Permissions []Permission `gorm:"many2many:role_permissions;"`
}



// struct for permissions
type Permission struct {
    gorm.Model
    Name string `gorm:"uniqueIndex"`
}

