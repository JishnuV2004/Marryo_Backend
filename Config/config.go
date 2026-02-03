package config

import (
	"log"
	models "marryo/Internal/Models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(){
	dsn := "host=localhost user=postgres password=jishnu@2004 dbname=marryo port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Profile{},
	)
	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	DB = db
	log.Println("Database connected")
}