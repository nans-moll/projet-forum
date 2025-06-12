package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	// Configuration de la base de données MySQL
	dbConfig := "root:@tcp(localhost:3306)/forum_db?parseTime=true"

	var err error
	DB, err = sql.Open("mysql", dbConfig)
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