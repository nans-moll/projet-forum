package models

import (
	"projet-forum/config"
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	ThreadID  int       `json:"thread_id"`
	AuthorID  int       `json:"author_id"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
}

// TableName retourne le nom de la table pour le modèle Message
func (Message) TableName() string {
	return "messages"
}

// CreateMessage crée un nouveau message
func (m *Message) CreateMessage() error {
	query := `
		INSERT INTO messages (thread_id, author_id, content, image_url, created_at, updated_at, likes, dislikes)
		VALUES (?, ?, ?, ?, ?, ?, 0, 0)`

	result, err := config.DB.Exec(
		query,
		m.ThreadID,
		m.AuthorID,
		m.Content,
		m.ImageURL,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	m.ID = int(id)
	return nil
}

// GetMessageByID récupère un message par son ID
func GetMessageByID(id int) (*Message, error) {
	message := &Message{}
	query := `SELECT * FROM messages WHERE id = ?`
	err := config.DB.QueryRow(query, id).Scan(
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
	return message, nil
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

// UpdateMessage met à jour un message
func (m *Message) UpdateMessage() error {
	query := `
		UPDATE messages 
		SET content = ?, image_url = ?, updated_at = ?
		WHERE id = ?`

	_, err := config.DB.Exec(
		query,
		m.Content,
		m.ImageURL,
		time.Now(),
		m.ID,
	)
	return err
}

// DeleteMessage supprime un message
func (m *Message) DeleteMessage() error {
	query := `DELETE FROM messages WHERE id = ?`
	_, err := config.DB.Exec(query, m.ID)
	return err
}

// LikeMessage ajoute un like à un message
func (m *Message) LikeMessage() error {
	query := `UPDATE messages SET likes = likes + 1 WHERE id = ?`
	_, err := config.DB.Exec(query, m.ID)
	if err != nil {
		return err
	}
	m.Likes++
	return nil
}

// DislikeMessage ajoute un dislike à un message
func (m *Message) DislikeMessage() error {
	query := `UPDATE messages SET dislikes = dislikes + 1 WHERE id = ?`
	_, err := config.DB.Exec(query, m.ID)
	if err != nil {
		return err
	}
	m.Dislikes++
	return nil
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
