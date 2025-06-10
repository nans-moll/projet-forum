package seeders

import (
	"database/sql"
	"fmt"
	"projet-forum/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initialise la connexion à la base de données pour les seeders
func InitDB() error {
	cfg := config.LoadConfig()
	dbConnString := cfg.GetDBConnString()

	var err error
	DB, err = sql.Open("postgres", dbConnString)
	if err != nil {
		return fmt.Errorf("erreur lors de la connexion à la base de données: %v", err)
	}

	// Test de la connexion
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("erreur lors du ping de la base de données: %v", err)
	}

	fmt.Println("Connexion à la base de données réussie pour les seeders!")
	return nil
}

// CloseDB ferme la connexion à la base de données
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
