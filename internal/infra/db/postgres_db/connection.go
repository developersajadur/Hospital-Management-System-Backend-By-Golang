package postgres_db

import (
	"fmt"
	"hospital_management_system/config"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.ENV.DBHost,
		config.ENV.DBUserName,
		config.ENV.DBPassword,
		config.ENV.DBName,
		config.ENV.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
		os.Exit(1)
	}

	DB = db
	log.Println("Database connected with GORM")
}
