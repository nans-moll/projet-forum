package config

import (
	"log"
	"os"
)

var AppConfig struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

func LoadConfig() {
	AppConfig.Port = "8080"
	AppConfig.DBHost = "localhost"
	AppConfig.DBPort = "3306"
	AppConfig.DBUser = "root"
	AppConfig.DBPassword = ""
	AppConfig.DBName = "projet_forum"
	AppConfig.JWTSecret = "your-secret-key-here"

	if host := os.Getenv("DB_HOST"); host != "" {
		AppConfig.DBHost = host
	}
	if user := os.Getenv("DB_USER"); user != "" {
		AppConfig.DBUser = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		AppConfig.DBPassword = password
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		AppConfig.DBName = dbName
	}

	log.Println("✅ Configuration chargée")
}
