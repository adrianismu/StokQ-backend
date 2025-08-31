package config

import (
	"log"
	"os"
	"stokq-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established successfully")

	// Run auto migration
	err = DB.AutoMigrate(
		&models.User{},
		&models.Product{},
	)
	if err != nil {
		log.Fatal("Failed to run database migration:", err)
	}

	log.Println("Database migration completed successfully")
}
