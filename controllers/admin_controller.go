package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"projet-forum/middleware"
	"strconv"
)

type AdminController struct {
	DB *sql.DB
}

// BanUser bannit un utilisateur
func (c *AdminController) BanUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	claims := middleware.GetUserFromContext(r)
	if claims == nil || claims.Role != "admin" {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	userID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid user ID",
		})
		return
	}

	// Mettre à jour le statut de l'utilisateur
	_, err = c.DB.Exec("UPDATE users SET is_banned = true WHERE id = ?", userID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error banning user",
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status:  "success",
		Message: "User banned successfully",
	})
}

// UnbanUser débannit un utilisateur
func (c *AdminController) UnbanUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	claims := middleware.GetUserFromContext(r)
	if claims == nil || claims.Role != "admin" {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	userID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid user ID",
		})
		return
	}

	// Mettre à jour le statut de l'utilisateur
	_, err = c.DB.Exec("UPDATE users SET is_banned = false WHERE id = ?", userID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error unbanning user",
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status:  "success",
		Message: "User unbanned successfully",
	})
}

// UpdateThreadStatus met à jour le statut d'un fil de discussion
func (c *AdminController) UpdateThreadStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	claims := middleware.GetUserFromContext(r)
	if claims == nil || claims.Role != "admin" {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	threadID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid thread ID",
		})
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Validation du statut
	switch req.Status {
	case "open", "closed", "archived":
		// Mettre à jour le statut
		_, err = c.DB.Exec("UPDATE threads SET status = ? WHERE id = ?", req.Status, threadID)
		if err != nil {
			middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
				Status:  "error",
				Message: "Error updating thread status",
			})
			return
		}
	default:
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid status",
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status:  "success",
		Message: "Thread status updated successfully",
	})
}

// GetAdminStats récupère les statistiques du forum
func (c *AdminController) GetAdminStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	claims := middleware.GetUserFromContext(r)
	if claims == nil || claims.Role != "admin" {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	// Récupérer les statistiques
	var stats struct {
		TotalUsers      int `json:"total_users"`
		TotalThreads    int `json:"total_threads"`
		TotalMessages   int `json:"total_messages"`
		BannedUsers     int `json:"banned_users"`
		ActiveUsers     int `json:"active_users"`
		OpenThreads     int `json:"open_threads"`
		ClosedThreads   int `json:"closed_threads"`
		ArchivedThreads int `json:"archived_threads"`
	}

	// Compter les utilisateurs
	err := c.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting user count",
		})
		return
	}

	// Compter les utilisateurs bannis
	err = c.DB.QueryRow("SELECT COUNT(*) FROM users WHERE is_banned = true").Scan(&stats.BannedUsers)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting banned user count",
		})
		return
	}

	// Compter les utilisateurs actifs (connectés dans les dernières 24h)
	err = c.DB.QueryRow("SELECT COUNT(*) FROM users WHERE last_connection > DATE_SUB(NOW(), INTERVAL 24 HOUR)").Scan(&stats.ActiveUsers)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting active user count",
		})
		return
	}

	// Compter les fils de discussion
	err = c.DB.QueryRow("SELECT COUNT(*) FROM threads").Scan(&stats.TotalThreads)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting thread count",
		})
		return
	}

	// Compter les fils de discussion par statut
	err = c.DB.QueryRow("SELECT COUNT(*) FROM threads WHERE status = 'open'").Scan(&stats.OpenThreads)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting open thread count",
		})
		return
	}

	err = c.DB.QueryRow("SELECT COUNT(*) FROM threads WHERE status = 'closed'").Scan(&stats.ClosedThreads)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting closed thread count",
		})
		return
	}

	err = c.DB.QueryRow("SELECT COUNT(*) FROM threads WHERE status = 'archived'").Scan(&stats.ArchivedThreads)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting archived thread count",
		})
		return
	}

	// Compter les messages
	err = c.DB.QueryRow("SELECT COUNT(*) FROM messages").Scan(&stats.TotalMessages)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting message count",
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   stats,
	})
}
