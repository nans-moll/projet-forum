package config

import "os"

// GetEnvOrDefault retourne la valeur de la variable d'environnement si elle existe,
// sinon retourne la valeur par défaut
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
