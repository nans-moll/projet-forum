package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("erreur MySQL: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("connexion MySQL échouée: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	log.Println("✅ Base de données MySQL connectée")
	return db, nil
}
