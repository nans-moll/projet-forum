package config

import (
	"fmt"
)

// Configuration globale de l'application
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	JWTSecret  string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:     GetEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     GetEnvOrDefault("DB_PORT", "3306"), // Port MySQL par défaut
		DBUser:     GetEnvOrDefault("DB_USER", "root"), // Utilisateur MySQL par défaut
		DBPassword: GetEnvOrDefault("DB_PASSWORD", ""),
		DBName:     GetEnvOrDefault("DB_NAME", "forum_db"),
		ServerPort: GetEnvOrDefault("SERVER_PORT", "8080"),
		JWTSecret:  GetEnvOrDefault("JWT_SECRET", "votre_clé_secrète_jwt"),
	}
}

func (c *Config) GetDBConnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName)
}
