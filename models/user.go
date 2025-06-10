package models

import (
	"projet-forum/config"
	"time"
)

type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	Role           string    `json:"role"`
	IsBanned       bool      `json:"is_banned"`
	CreatedAt      time.Time `json:"created_at"`
	LastConnection time.Time `json:"last_connection"`
	ProfilePicture string    `json:"profile_picture"`
	Biography      string    `json:"biography"`
	MessageCount   int       `json:"message_count"`
	ThreadCount    int       `json:"thread_count"`
}

// TableName retourne le nom de la table pour le modèle User
func (User) TableName() string {
	return "users"
}

// CreateUser crée un nouvel utilisateur dans la base de données
func (u *User) CreateUser() error {
	query := `
		INSERT INTO users (username, email, password_hash, role, created_at, last_connection)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	return config.DB.QueryRow(
		query,
		u.Username,
		u.Email,
		u.PasswordHash,
		"user", // Rôle par défaut
		time.Now(),
		time.Now(),
	).Scan(&u.ID)
}

// GetUserByID récupère un utilisateur par son ID
func GetUserByID(id int) (*User, error) {
	user := &User{}
	query := `SELECT * FROM users WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsBanned,
		&user.CreatedAt,
		&user.LastConnection,
		&user.ProfilePicture,
		&user.Biography,
		&user.MessageCount,
		&user.ThreadCount,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail récupère un utilisateur par son email
func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := `SELECT * FROM users WHERE email = $1`
	err := config.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsBanned,
		&user.CreatedAt,
		&user.LastConnection,
		&user.ProfilePicture,
		&user.Biography,
		&user.MessageCount,
		&user.ThreadCount,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUsername récupère un utilisateur par son nom d'utilisateur
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	query := `SELECT * FROM users WHERE username = $1`
	err := config.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsBanned,
		&user.CreatedAt,
		&user.LastConnection,
		&user.ProfilePicture,
		&user.Biography,
		&user.MessageCount,
		&user.ThreadCount,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateLastConnection met à jour la dernière connexion de l'utilisateur
func (u *User) UpdateLastConnection() error {
	query := `UPDATE users SET last_connection = $1 WHERE id = $2`
	_, err := config.DB.Exec(query, time.Now(), u.ID)
	return err
}

// UpdateUser met à jour les informations d'un utilisateur
func (u *User) UpdateUser() error {
	query := `
		UPDATE users 
		SET username = $1, email = $2, password_hash = $3, role = $4, is_banned = $5,
			profile_picture = $6, biography = $7, message_count = $8, thread_count = $9
		WHERE id = $10`

	_, err := config.DB.Exec(
		query,
		u.Username,
		u.Email,
		u.PasswordHash,
		u.Role,
		u.IsBanned,
		u.ProfilePicture,
		u.Biography,
		u.MessageCount,
		u.ThreadCount,
		u.ID,
	)
	return err
}
