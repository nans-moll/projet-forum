package models

import (
	"projet-forum/config"
	"time"
)

type ThreadStatus string

const (
	ThreadOpen     ThreadStatus = "open"
	ThreadClosed   ThreadStatus = "closed"
	ThreadArchived ThreadStatus = "archived"
)

type ThreadVisibility string

const (
	ThreadPublic  ThreadVisibility = "public"
	ThreadPrivate ThreadVisibility = "private"
)

type Thread struct {
	ID           int              `json:"id"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Tags         []string         `json:"tags"`
	Status       ThreadStatus     `json:"status"`
	Visibility   ThreadVisibility `json:"visibility"`
	AuthorID     int              `json:"author_id"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	MessageCount int              `json:"message_count"`
}

// TableName retourne le nom de la table pour le modèle Thread
func (Thread) TableName() string {
	return "threads"
}

// CreateThread crée un nouveau fil de discussion
func (t *Thread) CreateThread() error {
	query := `
		INSERT INTO threads (title, description, tags, status, visibility, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	return config.DB.QueryRow(
		query,
		t.Title,
		t.Description,
		t.Tags,
		t.Status,
		t.Visibility,
		t.AuthorID,
		time.Now(),
		time.Now(),
	).Scan(&t.ID)
}

// GetThreadByID récupère un fil de discussion par son ID
func GetThreadByID(id int) (*Thread, error) {
	thread := &Thread{}
	query := `SELECT * FROM threads WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(
		&thread.ID,
		&thread.Title,
		&thread.Description,
		&thread.Tags,
		&thread.Status,
		&thread.Visibility,
		&thread.AuthorID,
		&thread.CreatedAt,
		&thread.UpdatedAt,
		&thread.MessageCount,
	)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

// UpdateThread met à jour un fil de discussion
func (t *Thread) UpdateThread() error {
	query := `
		UPDATE threads 
		SET title = $1, description = $2, tags = $3, status = $4, visibility = $5, updated_at = $6
		WHERE id = $7`

	_, err := config.DB.Exec(
		query,
		t.Title,
		t.Description,
		t.Tags,
		t.Status,
		t.Visibility,
		time.Now(),
		t.ID,
	)
	return err
}

// DeleteThread supprime un fil de discussion
func (t *Thread) DeleteThread() error {
	query := `DELETE FROM threads WHERE id = $1`
	_, err := config.DB.Exec(query, t.ID)
	return err
}

// GetThreadsByTag récupère les fils de discussion par tag
func GetThreadsByTag(tag string, limit, offset int) ([]*Thread, error) {
	query := `
		SELECT * FROM threads 
		WHERE $1 = ANY(tags) AND status != 'archived'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := config.DB.Query(query, tag, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*Thread
	for rows.Next() {
		thread := &Thread{}
		err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Description,
			&thread.Tags,
			&thread.Status,
			&thread.Visibility,
			&thread.AuthorID,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.MessageCount,
		)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}
	return threads, nil
}

// GetThreadsByTitle récupère les fils de discussion par titre
func GetThreadsByTitle(title string, limit, offset int) ([]*Thread, error) {
	query := `
		SELECT * FROM threads 
		WHERE title ILIKE $1 AND status != 'archived'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := config.DB.Query(query, "%"+title+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*Thread
	for rows.Next() {
		thread := &Thread{}
		err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Description,
			&thread.Tags,
			&thread.Status,
			&thread.Visibility,
			&thread.AuthorID,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.MessageCount,
		)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}
	return threads, nil
}

// GetThreadsByAuthorID récupère les fils de discussion d'un auteur
func GetThreadsByAuthorID(authorID, limit, offset int) ([]*Thread, error) {
	query := `
		SELECT * FROM threads 
		WHERE author_id = $1 AND status != 'archived'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := config.DB.Query(query, authorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*Thread
	for rows.Next() {
		thread := &Thread{}
		err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Description,
			&thread.Tags,
			&thread.Status,
			&thread.Visibility,
			&thread.AuthorID,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.MessageCount,
		)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}
	return threads, nil
}
