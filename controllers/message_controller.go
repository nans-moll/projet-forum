package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"projet-forum/middleware"
	"projet-forum/models"
)

type MessageController struct {
	DB *sql.DB
}

// CreateMessage gère la création d'un nouveau message
func (c *MessageController) CreateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Vérifier l'authentification
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	var input struct {
		Content  string `json:"content"`
		ThreadID int64  `json:"thread_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Vérifier si le fil de discussion existe et est ouvert
	thread, err := models.GetThread(c.DB, input.ThreadID)
	if err != nil {
		middleware.SendJSON(w, http.StatusNotFound, middleware.Response{
			Status:  "error",
			Message: "Thread not found",
		})
		return
	}

	if thread.Status == "closed" {
		middleware.SendJSON(w, http.StatusForbidden, middleware.Response{
			Status:  "error",
			Message: "Thread is closed",
		})
		return
	}

	// Créer le message
	message, err := models.CreateMessage(c.DB, input.ThreadID, claims.UserID, input.Content, "")
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error creating message: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusCreated, middleware.Response{
		Status: "success",
		Data:   message,
	})
}

// GetMessage gère la récupération d'un message
func (c *MessageController) GetMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Récupérer l'ID du message
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid message ID",
		})
		return
	}

	// Récupérer le message
	message, err := models.GetMessage(c.DB, id)
	if err != nil {
		middleware.SendJSON(w, http.StatusNotFound, middleware.Response{
			Status:  "error",
			Message: "Message not found",
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   message,
	})
}

// ListMessages gère la liste des messages d'un fil de discussion
func (c *MessageController) ListMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Récupérer l'ID du fil de discussion
	threadIDStr := r.URL.Query().Get("thread_id")
	threadID, err := strconv.ParseInt(threadIDStr, 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid thread ID",
		})
		return
	}

	// Récupérer les paramètres de pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// Récupérer le paramètre de tri
	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "newest"
	}

	// Récupérer la liste des messages
	messages, err := models.ListMessages(c.DB, threadID, page, perPage, sortBy)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error listing messages: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   messages,
	})
}

// UpdateMessage gère la mise à jour d'un message
func (c *MessageController) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Vérifier l'authentification
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	// Récupérer l'ID du message
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid message ID",
		})
		return
	}

	// Récupérer le message
	message, err := models.GetMessage(c.DB, id)
	if err != nil {
		middleware.SendJSON(w, http.StatusNotFound, middleware.Response{
			Status:  "error",
			Message: "Message not found",
		})
		return
	}

	// Vérifier les permissions
	if message.AuthorID != claims.UserID && claims.Role != "admin" {
		middleware.SendJSON(w, http.StatusForbidden, middleware.Response{
			Status:  "error",
			Message: "Not authorized to update this message",
		})
		return
	}

	var input struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Mettre à jour le message
	message.Content = input.Content
	if err := message.UpdateMessage(c.DB); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error updating message: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   message,
	})
}

// DeleteMessage gère la suppression d'un message
func (c *MessageController) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Vérifier l'authentification
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	// Récupérer l'ID du message
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid message ID",
		})
		return
	}

	// Récupérer le message
	message, err := models.GetMessage(c.DB, id)
	if err != nil {
		middleware.SendJSON(w, http.StatusNotFound, middleware.Response{
			Status:  "error",
			Message: "Message not found",
		})
		return
	}

	// Vérifier les permissions
	if message.AuthorID != claims.UserID && claims.Role != "admin" {
		middleware.SendJSON(w, http.StatusForbidden, middleware.Response{
			Status:  "error",
			Message: "Not authorized to delete this message",
		})
		return
	}

	// Supprimer le message
	if err := models.DeleteMessage(c.DB, id); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error deleting message: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status:  "success",
		Message: "Message deleted successfully",
	})
}

// VoteMessage gère le vote sur un message
func (c *MessageController) VoteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Vérifier l'authentification
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	var input struct {
		MessageID int64  `json:"message_id"`
		VoteType  string `json:"vote_type"` // "like" ou "dislike"
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Vérifier le type de vote
	if input.VoteType != "like" && input.VoteType != "dislike" {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid vote type",
		})
		return
	}

	// Ajouter ou mettre à jour le vote
	if err := models.AddMessageReaction(c.DB, input.MessageID, claims.UserID, input.VoteType); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error voting on message: " + err.Error(),
		})
		return
	}

	// Mettre à jour les compteurs
	if err := models.UpdateMessageReactionCount(c.DB, input.MessageID); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error updating message reaction count: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status:  "success",
		Message: "Vote recorded successfully",
	})
}
