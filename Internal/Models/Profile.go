package models

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model

	// Relationship
	UserID uint `gorm:"uniqueIndex;not null"`

	// Step 1: Basic Details
	FullName     string    `gorm:"size:100"`
	DOB          time.Time
	MotherTongue string    `gorm:"size:50"`

	// Step 2: Personal / Religious
	Gender         string `gorm:"size:20"`
	Height         string `gorm:"size:30"`
	PhysicalStatus string `gorm:"size:30"`
	MaritalStatus  string `gorm:"size:30"`
	Religion       string `gorm:"size:50"`

	// Step 3: Location / Professional
	Country      string `gorm:"size:50"`
	Employment   string `gorm:"size:50"`
	Occupation   string `gorm:"size:100"`
	AnnualIncome int

	// Step 4: Horoscope
	Star  string `gorm:"size:30"`
	Raasi string `gorm:"size:30"`

	// Step 5: Education / Organization
	Education   string `gorm:"size:100"`
	College     string `gorm:"size:150"`
	Organization string `gorm:"size:150"`

	// Step 6: Lifestyle
	EatingHabit string `gorm:"size:50"`

	// Meta
	ProfileCompleted bool `gorm:"default:false"`
}