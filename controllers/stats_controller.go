package controllers

import (
	"database/sql"
	"net/http"

	"projet-forum/middleware"
)

type StatsController struct {
	DB *sql.DB
}

// GetStats récupère les statistiques globales du forum
func (c *StatsController) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Récupérer le nombre total d'utilisateurs
	var userCount int
	err := c.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting user count: " + err.Error(),
		})
		return
	}

	// Récupérer le nombre total de discussions
	var threadCount int
	err = c.DB.QueryRow("SELECT COUNT(*) FROM threads").Scan(&threadCount)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting thread count: " + err.Error(),
		})
		return
	}

	// Récupérer le nombre total de messages
	var messageCount int
	err = c.DB.QueryRow("SELECT COUNT(*) FROM messages").Scan(&messageCount)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting message count: " + err.Error(),
		})
		return
	}

	// Récupérer les discussions les plus récentes (limité à 5)
	rows, err := c.DB.Query(`
		SELECT t.id, t.title, t.description, t.created_at, u.username as author
		FROM threads t
		JOIN users u ON t.author_id = u.id
		WHERE t.status = 'open' AND t.visibility = 'public'
		ORDER BY t.created_at DESC
		LIMIT 5
	`)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting recent threads: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	var recentThreads []map[string]interface{}
	for rows.Next() {
		var thread struct {
			ID          int
			Title       string
			Description string
			CreatedAt   string
			Author      string
		}
		err := rows.Scan(&thread.ID, &thread.Title, &thread.Description, &thread.CreatedAt, &thread.Author)
		if err != nil {
			middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
				Status:  "error",
				Message: "Error scanning thread: " + err.Error(),
			})
			return
		}
		recentThreads = append(recentThreads, map[string]interface{}{
			"id":          thread.ID,
			"title":       thread.Title,
			"description": thread.Description,
			"created_at":  thread.CreatedAt,
			"author":      thread.Author,
		})
	}

	// Récupérer les utilisateurs les plus actifs (limité à 5)
	rows, err = c.DB.Query(`
		SELECT u.id, u.username, COUNT(m.id) as message_count
		FROM users u
		LEFT JOIN messages m ON u.id = m.author_id
		GROUP BY u.id, u.username
		ORDER BY message_count DESC
		LIMIT 5
	`)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting active users: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	var activeUsers []map[string]interface{}
	for rows.Next() {
		var user struct {
			ID           int
			Username     string
			MessageCount int
		}
		err := rows.Scan(&user.ID, &user.Username, &user.MessageCount)
		if err != nil {
			middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
				Status:  "error",
				Message: "Error scanning user: " + err.Error(),
			})
			return
		}
		activeUsers = append(activeUsers, map[string]interface{}{
			"id":            user.ID,
			"username":      user.Username,
			"message_count": user.MessageCount,
		})
	}

	// Construire la réponse
	response := map[string]interface{}{
		"user_count":     userCount,
		"thread_count":   threadCount,
		"message_count":  messageCount,
		"recent_threads": recentThreads,
		"active_users":   activeUsers,
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   response,
	})
}
