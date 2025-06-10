package controllers

import (
	"encoding/json"
	"net/http"
	"projet-forum/middleware"
	"projet-forum/models"
	"strconv"
)

type CreateThreadRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Visibility  string   `json:"visibility"`
}

// CreateThread crée un nouveau fil de discussion
func CreateThread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateThreadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation des données
	if len(req.Title) < 3 {
		http.Error(w, "Title must be at least 3 characters long", http.StatusBadRequest)
		return
	}

	if len(req.Description) < 10 {
		http.Error(w, "Description must be at least 10 characters long", http.StatusBadRequest)
		return
	}

	// Créer le fil de discussion
	thread := &models.Thread{
		Title:       req.Title,
		Description: req.Description,
		Tags:        req.Tags,
		Status:      models.ThreadOpen,
		Visibility:  models.ThreadVisibility(req.Visibility),
		AuthorID:    user.ID,
	}

	if err := thread.CreateThread(); err != nil {
		http.Error(w, "Error creating thread", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(thread)
}

// GetThread récupère un fil de discussion par son ID
func GetThread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	// Vérifier la visibilité du fil
	if thread.Visibility == models.ThreadPrivate {
		user, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// TODO: Vérifier si l'utilisateur est ami avec l'auteur du fil
		if user.ID != thread.AuthorID {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(thread)
}

// GetThreadsByTag récupère les fils de discussion par tag
func GetThreadsByTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tag := r.URL.Query().Get("tag")
	if tag == "" {
		http.Error(w, "Tag parameter is required", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Valeur par défaut
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	threads, err := models.GetThreadsByTag(tag, limit, offset)
	if err != nil {
		http.Error(w, "Error fetching threads", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(threads)
}

// UpdateThread met à jour un fil de discussion
func UpdateThread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
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

	// Vérifier si l'utilisateur est l'auteur ou un admin
	if user.ID != thread.AuthorID && user.Role != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var req CreateThreadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Mettre à jour le fil
	thread.Title = req.Title
	thread.Description = req.Description
	thread.Tags = req.Tags
	thread.Visibility = models.ThreadVisibility(req.Visibility)

	if err := thread.UpdateThread(); err != nil {
		http.Error(w, "Error updating thread", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(thread)
}

// DeleteThread supprime un fil de discussion
func DeleteThread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
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

	// Vérifier si l'utilisateur est l'auteur ou un admin
	if user.ID != thread.AuthorID && user.Role != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if err := thread.DeleteThread(); err != nil {
		http.Error(w, "Error deleting thread", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
