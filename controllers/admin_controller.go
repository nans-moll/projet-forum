package controllers

import (
	"encoding/json"
	"net/http"
	"projet-forum/middleware"
	"projet-forum/models"
	"strconv"
)

// BanUser bannit un utilisateur
func BanUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	targetUser, err := models.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Un admin ne peut pas bannir un autre admin
	if targetUser.Role == "admin" {
		http.Error(w, "Cannot ban an admin user", http.StatusForbidden)
		return
	}

	targetUser.IsBanned = true
	if err := targetUser.UpdateUser(); err != nil {
		http.Error(w, "Error banning user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(targetUser)
}

// UnbanUser débannit un utilisateur
func UnbanUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	targetUser, err := models.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	targetUser.IsBanned = false
	if err := targetUser.UpdateUser(); err != nil {
		http.Error(w, "Error unbanning user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(targetUser)
}

// UpdateThreadStatus met à jour le statut d'un fil de discussion
func UpdateThreadStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	threadID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	thread, err := models.GetThreadByID(threadID)
	if err != nil {
		http.Error(w, "Thread not found", http.StatusNotFound)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation du statut
	switch req.Status {
	case "open", "closed", "archived":
		thread.Status = models.ThreadStatus(req.Status)
	default:
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	if err := thread.UpdateThread(); err != nil {
		http.Error(w, "Error updating thread status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(thread)
}

// GetAdminStats récupère les statistiques du forum
func GetAdminStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// TODO: Implémenter la récupération des statistiques
	stats := struct {
		TotalUsers      int `json:"total_users"`
		TotalThreads    int `json:"total_threads"`
		TotalMessages   int `json:"total_messages"`
		BannedUsers     int `json:"banned_users"`
		ActiveUsers     int `json:"active_users"`
		OpenThreads     int `json:"open_threads"`
		ClosedThreads   int `json:"closed_threads"`
		ArchivedThreads int `json:"archived_threads"`
	}{
		// Ces valeurs devront être calculées à partir de la base de données
		TotalUsers:      0,
		TotalThreads:    0,
		TotalMessages:   0,
		BannedUsers:     0,
		ActiveUsers:     0,
		OpenThreads:     0,
		ClosedThreads:   0,
		ArchivedThreads: 0,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
