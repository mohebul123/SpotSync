package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mohebul123/SpotSync/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment")
	}

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL environment variable is required but not set")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("Failed to get generic database object:", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour * 1)

	DB = database
	err = database.AutoMigrate(&models.User{}, &models.ParkingZone{}, &models.Reservation{})
	if err != nil {
		log.Fatal("Failed to run database migration:", err)
	}
	log.Println("Database migration completed successfully.")
	log.Println("Database connection successfully established with pooling configuration.")
}
