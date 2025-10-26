package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	JWTSecret        string
	DBUserName       string
	DBPassword       string
	DBName           string
	DBHost           string
	DBPort           string
	DBSslMode        string
	EmailHost        string
	EmailPort        string
	EmailAppPassword string
	Email            string
}

var ENV *Config

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using system environment variables.")
		os.Exit(1)
	}
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("Required environment variable %s is not set or empty", key)
		os.Exit(1)
	}
	return value
}

func Init() {
	loadEnv()

	ENV = &Config{
		Port:             getEnv("PORT"),
		JWTSecret:        getEnv("JWT_SECRET"),
		DBUserName:       getEnv("DB_USER"),
		DBPassword:       getEnv("DB_PASSWORD"),
		DBName:           getEnv("DB_NAME"),
		DBHost:           getEnv("DB_HOST"),
		DBPort:           getEnv("DB_PORT"),
		DBSslMode:        getEnv("DB_SSLMODE"),
		EmailHost:        getEnv("EMAIL_HOST"),
		EmailPort:        getEnv("EMAIL_PORT"),
		EmailAppPassword: getEnv("EMAIL_APP_PASSWORD"),
		Email:            getEnv("EMAIL"),
	}
}
