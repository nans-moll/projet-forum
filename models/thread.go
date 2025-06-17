package models

import (
	"database/sql"
	"fmt"
	"log"
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
	ViewCount    int       `json:"view_count"`
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
	log.Printf("[DEBUG] GetThread (model) - Début de la fonction pour l'ID: %d", id)

	query := `
		SELECT t.id, t.title, t.description, t.tags, t.author_id, t.status, t.visibility, t.created_at, t.updated_at,
			   COUNT(DISTINCT m.id) as message_count,
			   u.username, u.email, u.role
		FROM threads t
		LEFT JOIN messages m ON t.id = m.thread_id
		LEFT JOIN users u ON t.author_id = u.id
		WHERE t.id = ?
		GROUP BY t.id, t.title, t.description, t.tags, t.author_id, t.status, t.visibility, t.created_at, t.updated_at,
				 u.username, u.email, u.role
	`
	log.Printf("[DEBUG] GetThread (model) - Requête SQL: %s", query)

	thread := &Thread{}
	var author User

	log.Printf("[DEBUG] GetThread (model) - Exécution de la requête")
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
		log.Printf("[DEBUG] GetThread (model) - Erreur lors de la récupération: %v", err)
		return nil, err
	}

	thread.Author = &author

	log.Printf("[DEBUG] GetThread (model) - Données récupérées:")
	log.Printf("[DEBUG] GetThread (model) - Thread ID: %d", thread.ID)
	log.Printf("[DEBUG] GetThread (model) - Thread Title: %s", thread.Title)
	log.Printf("[DEBUG] GetThread (model) - Thread Description: %s", thread.Description)
	log.Printf("[DEBUG] GetThread (model) - Thread Tags: %s", thread.Tags)
	log.Printf("[DEBUG] GetThread (model) - Thread AuthorID: %d", thread.AuthorID)
	log.Printf("[DEBUG] GetThread (model) - Thread MessageCount: %d", thread.MessageCount)
	log.Printf("[DEBUG] GetThread (model) - Author Username: %s", thread.Author.Username)
	log.Printf("[DEBUG] GetThread (model) - Author Email: %s", thread.Author.Email)
	log.Printf("[DEBUG] GetThread (model) - Author Role: %s", thread.Author.Role)

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

// ListThreads récupère une liste de fils de discussion
func ListThreads(db *sql.DB, page, limit int, status, visibility string) ([]map[string]interface{}, error) {
	offset := (page - 1) * limit
	query := `
		SELECT t.id, t.title, t.description, t.tags, t.status, t.visibility, t.created_at, t.message_count,
			   u.id as author_id, u.username as author_username, u.email as author_email, u.role as author_role
		FROM threads t
		JOIN users u ON t.author_id = u.id
		WHERE t.status = ? AND t.visibility = ?
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?
	`

	fmt.Printf("[DEBUG] ListThreads - Query: %s\n", query)
	fmt.Printf("[DEBUG] ListThreads - Params: status=%s, visibility=%s, limit=%d, offset=%d\n", status, visibility, limit, offset)

	rows, err := db.Query(query, status, visibility, limit, offset)
	if err != nil {
		fmt.Printf("[DEBUG] ListThreads - Query error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var threads []map[string]interface{}
	for rows.Next() {
		var thread struct {
			ID           int64  `json:"id"`
			Title        string `json:"title"`
			Description  string `json:"description"`
			Tags         string `json:"tags"`
			Status       string `json:"status"`
			Visibility   string `json:"visibility"`
			CreatedAt    string `json:"created_at"`
			MessageCount int    `json:"message_count"`
			AuthorID     int64  `json:"author_id"`
			AuthorName   string `json:"author_username"`
			AuthorEmail  string `json:"author_email"`
			AuthorRole   string `json:"author_role"`
		}
		if err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Description,
			&thread.Tags,
			&thread.Status,
			&thread.Visibility,
			&thread.CreatedAt,
			&thread.MessageCount,
			&thread.AuthorID,
			&thread.AuthorName,
			&thread.AuthorEmail,
			&thread.AuthorRole,
		); err != nil {
			fmt.Printf("[DEBUG] ListThreads - Scan error: %v\n", err)
			return nil, err
		}

		fmt.Printf("[DEBUG] ListThreads - Thread data: ID=%d, Title=%s, Author=%s\n",
			thread.ID, thread.Title, thread.AuthorName)

		threads = append(threads, map[string]interface{}{
			"id":            thread.ID,
			"title":         thread.Title,
			"description":   thread.Description,
			"tags":          thread.Tags,
			"status":        thread.Status,
			"visibility":    thread.Visibility,
			"created_at":    thread.CreatedAt,
			"message_count": thread.MessageCount,
			"author": map[string]interface{}{
				"id":       thread.AuthorID,
				"username": thread.AuthorName,
				"email":    thread.AuthorEmail,
				"role":     thread.AuthorRole,
			},
		})
	}

	fmt.Printf("[DEBUG] ListThreads - Found %d threads\n", len(threads))
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

// AdminUpdateThread permet à un admin de mettre à jour un fil de discussion
func AdminUpdateThread(db *sql.DB, threadID int64, status, visibility string) error {
	query := `
		UPDATE threads 
		SET status = ?, visibility = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := db.Exec(query, status, visibility, threadID)
	return err
}
