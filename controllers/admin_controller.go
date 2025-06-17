package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"projet-forum/middleware"
	"strconv"

	"github.com/gorilla/mux"
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

// ShowDashboard affiche le tableau de bord administrateur
func (c *AdminController) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/admin/dashboard.html")
}

// ListUsers liste tous les utilisateurs
func (c *AdminController) ListUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := c.DB.Query("SELECT id, username, email, role, is_banned, last_connection FROM users")
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting users",
		})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var user struct {
			ID             int64  `json:"id"`
			Username       string `json:"username"`
			Email          string `json:"email"`
			Role           string `json:"role"`
			Banned         bool   `json:"is_banned"`
			LastConnection string `json:"last_connection"`
		}
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Banned, &user.LastConnection); err != nil {
			middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
				Status:  "error",
				Message: "Error scanning users",
			})
			return
		}
		users = append(users, map[string]interface{}{
			"id":              user.ID,
			"username":        user.Username,
			"email":           user.Email,
			"role":            user.Role,
			"is_banned":       user.Banned,
			"last_connection": user.LastConnection,
		})
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   users,
	})
}

// ListThreads liste tous les fils de discussion
func (c *AdminController) ListThreads(w http.ResponseWriter, r *http.Request) {
	rows, err := c.DB.Query("SELECT id, title, status, created_at FROM threads")
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting threads",
		})
		return
	}
	defer rows.Close()

	var threads []map[string]interface{}
	for rows.Next() {
		var thread struct {
			ID        int64  `json:"id"`
			Title     string `json:"title"`
			Status    string `json:"status"`
			CreatedAt string `json:"created_at"`
		}
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.Status, &thread.CreatedAt); err != nil {
			middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
				Status:  "error",
				Message: "Error scanning threads",
			})
			return
		}
		threads = append(threads, map[string]interface{}{
			"id":         thread.ID,
			"title":      thread.Title,
			"status":     thread.Status,
			"created_at": thread.CreatedAt,
		})
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   threads,
	})
}

// DeleteThread supprime un fil de discussion
func (c *AdminController) DeleteThread(w http.ResponseWriter, r *http.Request) {
	threadID := mux.Vars(r)["id"]
	_, err := c.DB.Exec("DELETE FROM threads WHERE id = ?", threadID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error deleting thread",
		})
		return
	}
	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status:  "success",
		Message: "Thread deleted successfully",
	})
}

// DeleteMessage supprime un message
func (c *AdminController) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	messageID := mux.Vars(r)["id"]
	_, err := c.DB.Exec("DELETE FROM messages WHERE id = ?", messageID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error deleting message",
		})
		return
	}
	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status:  "success",
		Message: "Message deleted successfully",
	})
}
