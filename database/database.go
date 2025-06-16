package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB initialise la connexion à la base de données MySQL
func InitDB() {
	var err error
	// Connexion à MySQL via XAMPP avec la base de données forum_db
	DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/forum_db")
	if err != nil {
		log.Fatal(err)
	}

	// Vérifier la connexion
	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	// Créer les tables si elles n'existent pas
	createTables()
}

// createTables crée les tables nécessaires dans la base de données
func createTables() {
	// Table des utilisateurs
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`

	// Table des threads
	createThreadsTable := `
	CREATE TABLE IF NOT EXISTS threads (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		user_id INT NOT NULL,
		category VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	// Table des messages
	createMessagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INT AUTO_INCREMENT PRIMARY KEY,
		content TEXT NOT NULL,
		thread_id INT NOT NULL,
		user_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (thread_id) REFERENCES threads(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	// Table des likes
	createLikesTable := `
	CREATE TABLE IF NOT EXISTS likes (
		id INT AUTO_INCREMENT PRIMARY KEY,
		thread_id INT NOT NULL,
		user_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY unique_like (thread_id, user_id),
		FOREIGN KEY (thread_id) REFERENCES threads(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	// Exécuter les requêtes de création des tables
	_, err := DB.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(createThreadsTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(createMessagesTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(createLikesTable)
	if err != nil {
		log.Fatal(err)
	}
}

// CreateUser crée un nouvel utilisateur dans la base de données
func CreateUser(username, email, password string) error {
	query := `INSERT INTO users (username, email, mdp) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, username, email, password)
	return err
}

// GetUserByEmail récupère un utilisateur par son email
func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, email, mdp FROM users WHERE email = ?`
	err := DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// User représente un utilisateur dans la base de données
type User struct {
	ID       int64
	Username string
	Email    string
	Password string
}
