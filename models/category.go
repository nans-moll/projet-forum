package models

import (
	"database/sql"
	"time"
)

type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ThreadCount int       `json:"thread_count"`
}

// CreateCategory crée une nouvelle catégorie
func CreateCategory(db *sql.DB, name, description string) (*Category, error) {
	query := `
		INSERT INTO categories (name, description, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	result, err := db.Exec(query, name, description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetCategory(db, id)
}

// GetCategory récupère une catégorie par son ID
func GetCategory(db *sql.DB, id int64) (*Category, error) {
	query := `
		SELECT c.id, c.name, c.description, c.created_at, c.updated_at,
		       COUNT(t.id) as thread_count
		FROM categories c
		LEFT JOIN threads t ON c.id = t.category_id
		WHERE c.id = ?
		GROUP BY c.id
	`
	category := &Category{}
	err := db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.ThreadCount,
	)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// UpdateCategory met à jour une catégorie
func (c *Category) UpdateCategory(db *sql.DB) error {
	query := `
		UPDATE categories
		SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := db.Exec(query, c.Name, c.Description, c.ID)
	return err
}

// DeleteCategory supprime une catégorie
func DeleteCategory(db *sql.DB, id int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Supprimer d'abord les fils de discussion associés
	_, err = tx.Exec("DELETE FROM threads WHERE category_id = ?", id)
	if err != nil {
		return err
	}

	// Supprimer ensuite la catégorie
	_, err = tx.Exec("DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// ListCategories récupère la liste des catégories
func ListCategories(db *sql.DB) ([]*Category, error) {
	query := `
		SELECT c.id, c.name, c.description, c.created_at, c.updated_at,
		       COUNT(t.id) as thread_count
		FROM categories c
		LEFT JOIN threads t ON c.id = t.category_id
		GROUP BY c.id
		ORDER BY c.name ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		category := &Category{}
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.ThreadCount,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
