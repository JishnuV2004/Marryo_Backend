package config

import (
	"log"
	models "marryo/Internal/Models"
	"os"

	// "github.com/joho/godotenv"
	// "github.com/joho/godotenv"
	// "github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	// err := godotenv.Load("./.env")
	// 	if err != nil {
	//         log.Println("env file not found")
	//     }

	// root := os.Getenv("DB")
	// dns := "host=localhost user=postgres password=jishnu@2004 dbname=marryo port=5432 sslmode=disable"
	// 	dsn:=fmt.Sprintf(
	//   "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
	//   os.Getenv("DB_HOST"),
	//   os.Getenv("DB_USER"),
	//   os.Getenv("DB_PASSWORD"),
	//   os.Getenv("DB_NAME"),
	//   os.Getenv("DB_PORT"),
	//   os.Getenv("DB_SSLMODE"),
	//  )
	//  fmt.Println("DSN:", dsn)
	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Img{},
		&models.Interest{},
		&models.Match{},
		&models.Message{},
		&models.DeviceToken{},
		&models.Permission{},
		&models.Role{},
	)
	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	DB = db
	log.Println("Database connected")
}
