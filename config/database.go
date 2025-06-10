package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	// Configuration de la base de données PostgreSQL
	dbConfig := "host=localhost port=5432 user=postgres password=postgres dbname=forum_db sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", dbConfig)
	if err != nil {
		return fmt.Errorf("erreur lors de la connexion à la base de données: %v", err)
	}

	// Test de la connexion
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("erreur lors du ping de la base de données: %v", err)
	}

	fmt.Println("Connexion à la base de données réussie!")
	return nil
}
