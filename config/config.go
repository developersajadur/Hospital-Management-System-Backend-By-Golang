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
	RabbitMqUrl      string
	SendGridApiKey   string
	CloudinaryCloudName string
	CloudinaryApiKey string
	CloudinaryApiSecret string
	SSLStoreID       string
	SSLStorePassword string
	SSlSandbox      string
	BaseURL          string
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
		RabbitMqUrl:      getEnv("RABBITMQ_URL"),
		SendGridApiKey:   getEnv("SENDGRID_API_KEY"),
		CloudinaryCloudName: getEnv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryApiKey: getEnv("CLOUDINARY_API_KEY"),
		CloudinaryApiSecret: getEnv("CLOUDINARY_API_SECRET"),
		SSLStoreID:       getEnv("SSL_STORE_ID"),
		SSLStorePassword: getEnv("SSL_STORE_PASSWORD"),
		SSlSandbox:      getEnv("SSL_SANDBOX"),
		BaseURL:          getEnv("BASE_URL"),

	}
}
