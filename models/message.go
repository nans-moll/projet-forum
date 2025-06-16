package models

import (
	"database/sql"
	"projet-forum/config"
	"time"
)

type Message struct {
	ID        int64     `json:"id"`
	ThreadID  int64     `json:"thread_id"`
	AuthorID  int64     `json:"author_id"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	Author    *User     `json:"author,omitempty"`
}

// TableName retourne le nom de la table pour le modèle Message
func (Message) TableName() string {
	return "messages"
}

// GetMessage récupère un message par son ID
func GetMessage(db *sql.DB, id int64) (*Message, error) {
	query := `
		SELECT m.id, m.thread_id, m.author_id, m.content, m.image_url, m.created_at, m.updated_at, m.likes, m.dislikes,
		       u.username, u.email, u.role
		FROM messages m
		LEFT JOIN users u ON m.author_id = u.id
		WHERE m.id = ?
	`
	message := &Message{}
	var author User
	err := db.QueryRow(query, id).Scan(
		&message.ID,
		&message.ThreadID,
		&message.AuthorID,
		&message.Content,
		&message.ImageURL,
		&message.CreatedAt,
		&message.UpdatedAt,
		&message.Likes,
		&message.Dislikes,
		&author.Username,
		&author.Email,
		&author.Role,
	)
	if err != nil {
		return nil, err
	}
	message.Author = &author
	return message, nil
}

// UpdateMessage met à jour un message
func (m *Message) UpdateMessage(db *sql.DB) error {
	query := `
		UPDATE messages
		SET content = ?, image_url = ?, likes = ?, dislikes = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := db.Exec(query, m.Content, m.ImageURL, m.Likes, m.Dislikes, m.ID)
	return err
}

// CreateMessage crée un nouveau message
func CreateMessage(db *sql.DB, threadID, authorID int64, content, imageURL string) (*Message, error) {
	query := `
		INSERT INTO messages (thread_id, author_id, content, image_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	result, err := db.Exec(query, threadID, authorID, content, imageURL)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Mettre à jour le compteur de messages du fil de discussion
	_, err = db.Exec("UPDATE threads SET message_count = message_count + 1 WHERE id = ?", threadID)
	if err != nil {
		return nil, err
	}

	// Mettre à jour le compteur de messages de l'auteur
	_, err = db.Exec("UPDATE users SET message_count = message_count + 1 WHERE id = ?", authorID)
	if err != nil {
		return nil, err
	}

	return GetMessage(db, id)
}

// DeleteMessage supprime un message
func DeleteMessage(db *sql.DB, id int64) error {
	query := `DELETE FROM messages WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// ListMessages récupère la liste des messages d'un fil de discussion
func ListMessages(db *sql.DB, threadID int64, page, perPage int, sortBy string) ([]*Message, error) {
	offset := (page - 1) * perPage
	var orderBy string

	switch sortBy {
	case "likes":
		orderBy = "m.likes - m.dislikes DESC"
	case "oldest":
		orderBy = "m.created_at ASC"
	default:
		orderBy = "m.created_at DESC"
	}

	query := `
		SELECT m.id, m.thread_id, m.author_id, m.content, m.image_url, m.created_at, m.updated_at, m.likes, m.dislikes,
		       u.username, u.email, u.role
		FROM messages m
		LEFT JOIN users u ON m.author_id = u.id
		WHERE m.thread_id = ?
		ORDER BY ` + orderBy + `
		LIMIT ? OFFSET ?
	`

	rows, err := db.Query(query, threadID, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		message := &Message{}
		var author User
		err := rows.Scan(
			&message.ID,
			&message.ThreadID,
			&message.AuthorID,
			&message.Content,
			&message.ImageURL,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.Likes,
			&message.Dislikes,
			&author.Username,
			&author.Email,
			&author.Role,
		)
		if err != nil {
			return nil, err
		}
		message.Author = &author
		messages = append(messages, message)
	}

	return messages, nil
}

// GetMessagesByThreadID récupère tous les messages d'un fil de discussion
func GetMessagesByThreadID(threadID int, limit, offset int, sortBy string) ([]*Message, error) {
	var orderBy string
	switch sortBy {
	case "popularity":
		orderBy = "likes - dislikes DESC"
	case "oldest":
		orderBy = "created_at ASC"
	default: // "newest"
		orderBy = "created_at DESC"
	}

	query := `
		SELECT * FROM messages 
		WHERE thread_id = ?
		ORDER BY ` + orderBy + `
		LIMIT ? OFFSET ?`

	rows, err := config.DB.Query(query, threadID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		message := &Message{}
		err := rows.Scan(
			&message.ID,
			&message.ThreadID,
			&message.AuthorID,
			&message.Content,
			&message.ImageURL,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.Likes,
			&message.Dislikes,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// GetMessagesByAuthorID récupère les messages d'un auteur
func GetMessagesByAuthorID(authorID, limit, offset int) ([]*Message, error) {
	query := `
		SELECT * FROM messages 
		WHERE author_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := config.DB.Query(query, authorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		message := &Message{}
		err := rows.Scan(
			&message.ID,
			&message.ThreadID,
			&message.AuthorID,
			&message.Content,
			&message.ImageURL,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.Likes,
			&message.Dislikes,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
