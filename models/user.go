package models

import (
	"database/sql"
	"time"

	"crypto/sha512"
	"encoding/hex"

	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	Role           string    `json:"role"`
	Banned         bool      `json:"banned"`
	ThreadCount    int       `json:"thread_count"`
	MessageCount   int       `json:"message_count"`
	LastConnection time.Time `json:"last_connection"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	ProfilePicture string    `json:"profile_picture,omitempty"`
	Biography      string    `json:"biography,omitempty"`
}

// TableName retourne le nom de la table pour le modèle User
func (User) TableName() string {
	return "users"
}

// CreateUser crée un nouvel utilisateur
func CreateUser(db *sql.DB, username, email, password, role string) (*User, error) {
	query := `
		INSERT INTO users (username, email, password_hash, role, is_banned, thread_count, message_count, last_connection)
		VALUES (?, ?, ?, ?, ?, 0, 0, CURRENT_TIMESTAMP)
	`
	result, err := db.Exec(query, username, email, password, role, false)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetUserByID(db, id)
}

// GetUserByID récupère un utilisateur par son ID
func GetUserByID(db *sql.DB, id int64) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, role, is_banned, thread_count, message_count, last_connection, created_at
		FROM users
		WHERE id = ?
	`
	user := &User{}
	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Banned,
		&user.ThreadCount,
		&user.MessageCount,
		&user.LastConnection,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail récupère un utilisateur par son email
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, role, is_banned, thread_count, message_count, last_connection, created_at
		FROM users WHERE email = ?
	`
	user := &User{}
	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Banned,
		&user.ThreadCount,
		&user.MessageCount,
		&user.LastConnection,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUsername récupère un utilisateur par son nom d'utilisateur
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, role, is_banned, thread_count, message_count, last_connection, created_at
		FROM users
		WHERE username = ?
	`
	user := &User{}
	err := db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Banned,
		&user.ThreadCount,
		&user.MessageCount,
		&user.LastConnection,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser met à jour un utilisateur
func (u *User) UpdateUser(db *sql.DB) error {
	query := `
		UPDATE users
		SET username = ?, email = ?, password_hash = ?, role = ?, is_banned = ?
		WHERE id = ?
	`
	_, err := db.Exec(query, u.Username, u.Email, u.PasswordHash, u.Role, u.Banned, u.ID)
	return err
}

// DeleteUser supprime un utilisateur
func DeleteUser(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

// ListUsers récupère une liste d'utilisateurs
func ListUsers(db *sql.DB, page, limit int) ([]map[string]interface{}, error) {
	offset := (page - 1) * limit
	query := `
		SELECT id, username, email, role, thread_count, message_count, last_connection
		FROM users
		ORDER BY message_count DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var user struct {
			ID             int64  `json:"id"`
			Username       string `json:"username"`
			Email          string `json:"email"`
			Role           string `json:"role"`
			ThreadCount    int    `json:"thread_count"`
			MessageCount   int    `json:"message_count"`
			LastConnection string `json:"last_connection"`
		}
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.ThreadCount, &user.MessageCount, &user.LastConnection); err != nil {
			return nil, err
		}
		users = append(users, map[string]interface{}{
			"id":              user.ID,
			"username":        user.Username,
			"email":           user.Email,
			"role":            user.Role,
			"thread_count":    user.ThreadCount,
			"message_count":   user.MessageCount,
			"last_connection": user.LastConnection,
		})
	}
	return users, nil
}

// CheckPassword vérifie si le mot de passe est correct
func (u *User) CheckPassword(password string) bool {
	// D'abord essayer avec bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err == nil {
		return true
	}

	// Si bcrypt échoue, essayer avec SHA512
	hashedPassword := sha512.Sum512([]byte(password))
	hashedPasswordStr := hex.EncodeToString(hashedPassword[:])
	return hashedPasswordStr == u.PasswordHash
}

// SetPassword définit un nouveau mot de passe
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

// UpdateLastConnection met à jour la dernière connexion
func (u *User) UpdateLastConnection(db *sql.DB) error {
	_, err := db.Exec("UPDATE users SET last_connection = CURRENT_TIMESTAMP WHERE id = ?", u.ID)
	return err
}

// UpdateLastLogin met à jour la dernière connexion
func (u *User) UpdateLastLogin(db *sql.DB) error {
	_, err := db.Exec("UPDATE users SET last_login = NOW() WHERE id = ?", u.ID)
	return err
}

// BanUser bannit un utilisateur
func BanUser(db *sql.DB, userID int64) error {
	_, err := db.Exec("UPDATE users SET banned = true WHERE id = ?", userID)
	return err
}

// UnbanUser débannit un utilisateur
func UnbanUser(db *sql.DB, userID int64) error {
	_, err := db.Exec("UPDATE users SET banned = false WHERE id = ?", userID)
	return err
}

// ValidatePassword vérifie si le mot de passe est correct
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// IsAdmin vérifie si l'utilisateur est un administrateur
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsBanned vérifie si l'utilisateur est banni
func (u *User) IsBanned() bool {
	return u.Banned
}

// GetUserStats récupère les statistiques d'un utilisateur
func GetUserStats(db *sql.DB, userID int64) (map[string]int, error) {
	stats := make(map[string]int)
	var count int

	// Nombre de fils de discussion
	err := db.QueryRow("SELECT COUNT(*) FROM threads WHERE author_id = ?", userID).Scan(&count)
	if err != nil {
		return nil, err
	}
	stats["threads"] = count

	// Nombre de messages
	err = db.QueryRow("SELECT COUNT(*) FROM messages WHERE author_id = ?", userID).Scan(&count)
	if err != nil {
		return nil, err
	}
	stats["messages"] = count

	// Nombre total de likes reçus
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM message_reactions mr
		JOIN messages m ON mr.message_id = m.id
		WHERE m.author_id = ? AND mr.reaction_type = 'like'
	`, userID).Scan(&count)
	if err != nil {
		return nil, err
	}
	stats["likes"] = count

	return stats, nil
}

// Update met à jour les informations de l'utilisateur
func (u *User) Update(db *sql.DB) error {
	query := `
		UPDATE users 
		SET username = ?, email = ?
		WHERE id = ?
	`
	_, err := db.Exec(query, u.Username, u.Email, u.ID)
	return err
}

// GetUserMessages récupère tous les messages d'un utilisateur
func GetUserMessages(db *sql.DB, userID int64) ([]map[string]interface{}, error) {
	query := `
		SELECT m.id, m.content, m.created_at, t.title as thread_title
		FROM messages m
		JOIN threads t ON m.thread_id = t.id
		WHERE m.user_id = ?
		ORDER BY m.created_at DESC
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		var msg struct {
			ID          int64  `json:"id"`
			Content     string `json:"content"`
			CreatedAt   string `json:"created_at"`
			ThreadTitle string `json:"thread_title"`
		}
		if err := rows.Scan(&msg.ID, &msg.Content, &msg.CreatedAt, &msg.ThreadTitle); err != nil {
			return nil, err
		}
		messages = append(messages, map[string]interface{}{
			"id":           msg.ID,
			"content":      msg.Content,
			"created_at":   msg.CreatedAt,
			"thread_title": msg.ThreadTitle,
		})
	}
	return messages, nil
}

// GetUserThreads récupère tous les fils de discussion d'un utilisateur
func GetUserThreads(db *sql.DB, userID int64) ([]map[string]interface{}, error) {
	query := `
		SELECT id, title, description, status, created_at
		FROM threads
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []map[string]interface{}
	for rows.Next() {
		var thread struct {
			ID          int64  `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Status      string `json:"status"`
			CreatedAt   string `json:"created_at"`
		}
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.Description, &thread.Status, &thread.CreatedAt); err != nil {
			return nil, err
		}
		threads = append(threads, map[string]interface{}{
			"id":          thread.ID,
			"title":       thread.Title,
			"description": thread.Description,
			"status":      thread.Status,
			"created_at":  thread.CreatedAt,
		})
	}
	return threads, nil
}

// AuthenticateUser vérifie les identifiants d'un utilisateur
func AuthenticateUser(db *sql.DB, username, password string) (*User, error) {
	var user User
	query := `
		SELECT id, username, email, password_hash, role, is_banned, created_at, last_connection,
			   profile_picture, biography, message_count, thread_count
		FROM users
		WHERE (username = ? OR email = ?) AND is_banned = false
	`
	err := db.QueryRow(query, username, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.Banned, &user.CreatedAt, &user.LastConnection,
		&user.ProfilePicture, &user.Biography, &user.MessageCount, &user.ThreadCount,
	)
	if err != nil {
		return nil, err
	}

	// Vérifier le mot de passe
	if !user.CheckPassword(password) {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}

// UpdateLastConnection met à jour la dernière connexion d'un utilisateur
func UpdateLastConnection(db *sql.DB, userID int64) error {
	query := `
		UPDATE users
		SET last_connection = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := db.Exec(query, userID)
	return err
}
