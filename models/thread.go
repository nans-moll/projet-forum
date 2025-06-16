package models

import (
	"database/sql"
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
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Tags         string    `json:"tags"`
	AuthorID     int64     `json:"author_id"`
	Status       string    `json:"status"`
	Visibility   string    `json:"visibility"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	MessageCount int       `json:"message_count"`
	Author       *User     `json:"author,omitempty"`
}

// TableName retourne le nom de la table pour le modèle Thread
func (Thread) TableName() string {
	return "threads"
}

// CreateThread crée un nouveau fil de discussion
func CreateThread(db *sql.DB, title, description string, authorID int64, tags []string) (*Thread, error) {
	tagsStr := ""
	if len(tags) > 0 {
		for i, tag := range tags {
			if i > 0 {
				tagsStr += ","
			}
			tagsStr += tag
		}
	}

	query := `
		INSERT INTO threads (title, description, tags, author_id, status, visibility, created_at, updated_at)
		VALUES (?, ?, ?, ?, 'open', 'public', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	result, err := db.Exec(query, title, description, tagsStr, authorID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetThread(db, id)
}

// GetThread récupère un fil de discussion par son ID
func GetThread(db *sql.DB, id int64) (*Thread, error) {
	query := `
		SELECT t.id, t.title, t.description, t.tags, t.author_id, t.status, t.visibility, t.created_at, t.updated_at,
		       COUNT(m.id) as message_count,
		       u.username, u.email, u.role
		FROM threads t
		LEFT JOIN messages m ON t.id = m.thread_id
		LEFT JOIN users u ON t.author_id = u.id
		WHERE t.id = ?
		GROUP BY t.id
	`
	thread := &Thread{}
	var author User
	err := db.QueryRow(query, id).Scan(
		&thread.ID,
		&thread.Title,
		&thread.Description,
		&thread.Tags,
		&thread.AuthorID,
		&thread.Status,
		&thread.Visibility,
		&thread.CreatedAt,
		&thread.UpdatedAt,
		&thread.MessageCount,
		&author.Username,
		&author.Email,
		&author.Role,
	)
	if err != nil {
		return nil, err
	}
	thread.Author = &author
	return thread, nil
}

// UpdateThread met à jour un fil de discussion
func (t *Thread) UpdateThread(db *sql.DB) error {
	query := `
		UPDATE threads
		SET title = ?, description = ?, tags = ?, status = ?, visibility = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := db.Exec(query, t.Title, t.Description, t.Tags, t.Status, t.Visibility, t.ID)
	return err
}

// DeleteThread supprime un fil de discussion
func DeleteThread(db *sql.DB, id int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Supprimer d'abord les messages associés
	_, err = tx.Exec("DELETE FROM messages WHERE thread_id = ?", id)
	if err != nil {
		return err
	}

	// Supprimer ensuite le fil de discussion
	_, err = tx.Exec("DELETE FROM threads WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// ListThreads récupère la liste des fils de discussion
func ListThreads(db *sql.DB, page, perPage int, status string, tag string) ([]*Thread, error) {
	offset := (page - 1) * perPage
	var args []interface{}

	baseQuery := `
		SELECT t.id, t.title, t.description, t.tags, t.author_id, t.status, t.visibility, t.created_at, t.updated_at,
		       COUNT(m.id) as message_count,
		       u.username, u.email, u.role
		FROM threads t
		LEFT JOIN messages m ON t.id = m.thread_id
		LEFT JOIN users u ON t.author_id = u.id
	`

	whereClause := " WHERE t.status != 'archived'"
	if status != "" {
		whereClause += " AND t.status = ?"
		args = append(args, status)
	}
	if tag != "" {
		whereClause += " AND FIND_IN_SET(?, t.tags)"
		args = append(args, tag)
	}

	baseQuery += whereClause + " GROUP BY t.id ORDER BY t.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, perPage, offset)

	rows, err := db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*Thread
	for rows.Next() {
		thread := &Thread{}
		var author User
		err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Description,
			&thread.Tags,
			&thread.AuthorID,
			&thread.Status,
			&thread.Visibility,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.MessageCount,
			&author.Username,
			&author.Email,
			&author.Role,
		)
		if err != nil {
			return nil, err
		}
		thread.Author = &author
		threads = append(threads, thread)
	}
	return threads, nil
}

// UpdateThreadStatus met à jour le statut d'un fil de discussion
func UpdateThreadStatus(db *sql.DB, id int64, status string) error {
	_, err := db.Exec("UPDATE threads SET status = ? WHERE id = ?", status, id)
	return err
}

// SearchThreads recherche des fils de discussion
func SearchThreads(db *sql.DB, query string, page, perPage int) ([]*Thread, error) {
	offset := (page - 1) * perPage
	searchPattern := "%" + query + "%"

	sqlQuery := `
		SELECT t.id, t.title, t.description, t.tags, t.author_id, t.status, t.visibility, t.created_at, t.updated_at,
		       COUNT(m.id) as message_count,
		       u.username, u.email, u.role
		FROM threads t
		LEFT JOIN messages m ON t.id = m.thread_id
		LEFT JOIN users u ON t.author_id = u.id
		WHERE (t.title LIKE ? OR t.description LIKE ? OR FIND_IN_SET(?, t.tags))
		AND t.status != 'archived'
		GROUP BY t.id
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := db.Query(sqlQuery, searchPattern, searchPattern, query, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*Thread
	for rows.Next() {
		thread := &Thread{}
		var author User
		err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Description,
			&thread.Tags,
			&thread.AuthorID,
			&thread.Status,
			&thread.Visibility,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.MessageCount,
			&author.Username,
			&author.Email,
			&author.Role,
		)
		if err != nil {
			return nil, err
		}
		thread.Author = &author
		threads = append(threads, thread)
	}

	return threads, nil
}

// GetThreadsByTag récupère les fils de discussion par tag
func GetThreadsByTag(tag string, limit, offset int) ([]*Thread, error) {
	query := `
		SELECT * FROM threads 
		WHERE FIND_IN_SET(?, tags) AND status != 'archived'
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

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
			&thread.AuthorID,
			&thread.Status,
			&thread.Visibility,
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
		WHERE title LIKE ? AND status != 'archived'
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

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
			&thread.AuthorID,
			&thread.Status,
			&thread.Visibility,
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
		WHERE author_id = ? AND status != 'archived'
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

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
			&thread.AuthorID,
			&thread.Status,
			&thread.Visibility,
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
