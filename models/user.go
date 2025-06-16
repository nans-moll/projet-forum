package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Password       string    `json:"-"`
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
		INSERT INTO users (username, email, password, role, banned, thread_count, message_count, last_connection)
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
		SELECT id, username, email, password, role, is_banned, thread_count, message_count, last_connection, created_at, updated_at
		FROM users
		WHERE id = ?
	`
	user := &User{}
	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Banned,
		&user.ThreadCount,
		&user.MessageCount,
		&user.LastConnection,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail récupère un utilisateur par son email
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, role, is_banned, thread_count, message_count, last_connection, created_at, updated_at
		FROM users WHERE email = ?
	`
	user := &User{}
	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Banned,
		&user.ThreadCount,
		&user.MessageCount,
		&user.LastConnection,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUsername récupère un utilisateur par son nom d'utilisateur
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, role, is_banned, thread_count, message_count, last_connection, created_at, updated_at
		FROM users
		WHERE username = ?
	`
	user := &User{}
	err := db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Banned,
		&user.ThreadCount,
		&user.MessageCount,
		&user.LastConnection,
		&user.CreatedAt,
		&user.UpdatedAt,
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
		SET username = ?, email = ?, password = ?, role = ?, is_banned = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := db.Exec(query, u.Username, u.Email, u.Password, u.Role, u.Banned, u.ID)
	return err
}

// DeleteUser supprime un utilisateur
func DeleteUser(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

// ListUsers récupère la liste des utilisateurs
func ListUsers(db *sql.DB, page, perPage int) ([]*User, error) {
	offset := (page - 1) * perPage
	query := `
		SELECT id, username, email, role, banned, thread_count, message_count, last_connection, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.Banned,
			&user.ThreadCount,
			&user.MessageCount,
			&user.LastConnection,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// CheckPassword vérifie si le mot de passe est correct
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// SetPassword définit un nouveau mot de passe
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
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
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
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
