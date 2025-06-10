package controllers

import (
	"encoding/json"
	"net/http"
	"projet-forum/middleware"
	"projet-forum/models"
	"strconv"
)

type CreateMessageRequest struct {
	Content  string `json:"content"`
	ImageURL string `json:"image_url,omitempty"`
}

// CreateMessage crée un nouveau message dans un fil de discussion
func CreateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	threadID, err := strconv.Atoi(r.URL.Query().Get("thread_id"))
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	// Vérifier si le fil existe et est ouvert
	thread, err := models.GetThreadByID(threadID)
	if err != nil {
		http.Error(w, "Thread not found", http.StatusNotFound)
		return
	}

	if thread.Status != models.ThreadOpen {
		http.Error(w, "Thread is closed or archived", http.StatusForbidden)
		return
	}

	var req CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation du contenu
	if len(req.Content) < 1 {
		http.Error(w, "Message content cannot be empty", http.StatusBadRequest)
		return
	}

	// Créer le message
	message := &models.Message{
		ThreadID: threadID,
		AuthorID: user.ID,
		Content:  req.Content,
		ImageURL: req.ImageURL,
	}

	if err := message.CreateMessage(); err != nil {
		http.Error(w, "Error creating message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

// GetMessages récupère les messages d'un fil de discussion
func GetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	threadID, err := strconv.Atoi(r.URL.Query().Get("thread_id"))
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	// Vérifier si le fil existe
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

		if user.ID != thread.AuthorID {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "newest"
	}

	messages, err := models.GetMessagesByThreadID(threadID, limit, offset, sortBy)
	if err != nil {
		http.Error(w, "Error fetching messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	message, err := models.GetMessageByID(messageID)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	if user.ID != message.AuthorID && user.Role != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var req CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	message.Content = req.Content
	message.ImageURL = req.ImageURL

	if err := message.UpdateMessage(); err != nil {
		http.Error(w, "Error updating message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	message, err := models.GetMessageByID(messageID)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	if user.ID != message.AuthorID && user.Role != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if err := message.DeleteMessage(); err != nil {
		http.Error(w, "Error deleting message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func LikeMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	message, err := models.GetMessageByID(messageID)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	if err := message.LikeMessage(); err != nil {
		http.Error(w, "Error liking message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func DislikeMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	message, err := models.GetMessageByID(messageID)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	if err := message.DislikeMessage(); err != nil {
		http.Error(w, "Error disliking message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}
