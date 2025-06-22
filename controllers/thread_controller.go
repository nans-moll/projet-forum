package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"projet-forum/middleware"
	"projet-forum/models"

	"github.com/gorilla/mux"
)

// ThreadController gère les opérations liées aux fils de discussion
type ThreadController struct {
	DB        *sql.DB
	Templates *template.Template
}

// NewThreadController crée une nouvelle instance de ThreadController
func NewThreadController(db *sql.DB, templates *template.Template) *ThreadController {
	return &ThreadController{
		DB:        db,
		Templates: templates,
	}
}

// ShowThreadsPage affiche la page des fils de discussion
func (c *ThreadController) ShowThreadsPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/threads/index.html")
}

// ShowThreadPage affiche la page d'un fil de discussion
func (c *ThreadController) ShowThreadPage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[DEBUG] ShowThreadPage - URL Path: %s\n", r.URL.Path)

	// Extraire l'ID de la discussion de l'URL
	parts := strings.Split(r.URL.Path, "/")
	fmt.Printf("[DEBUG] ShowThreadPage - URL Parts: %v\n", parts)

	if len(parts) < 4 {
		fmt.Printf("[DEBUG] ShowThreadPage - Invalid URL format\n")
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}
	threadID := parts[3]
	fmt.Printf("[DEBUG] ShowThreadPage - Thread ID: %s\n", threadID)

	// Vérifier que l'ID est un nombre valide
	_, err := strconv.ParseInt(threadID, 10, 64)
	if err != nil {
		fmt.Printf("[DEBUG] ShowThreadPage - Invalid thread ID format: %v\n", err)
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	fmt.Printf("[DEBUG] ShowThreadPage - Serving file: templates/threads/show.html\n")
	http.ServeFile(w, r, "templates/threads/show.html")
}

// ShowSearchPage affiche la page de recherche
func (c *ThreadController) ShowSearchPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/search.html")
}

// CreateThread gère la création d'un nouveau fil de discussion
func (c *ThreadController) CreateThread(w http.ResponseWriter, r *http.Request) {
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
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Créer le fil de discussion
	thread, err := models.CreateThread(c.DB, input.Title, input.Description, claims.UserID, input.Tags)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error creating thread: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusCreated, middleware.Response{
		Status: "success",
		Data:   thread,
	})
}

// GetThread affiche une discussion et ses messages
func (c *ThreadController) GetThread(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DEBUG] GetThread - Début de la fonction")

	if r.Method != http.MethodGet {
		log.Printf("[DEBUG] GetThread - Méthode non autorisée: %s", r.Method)
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Récupérer l'ID du fil de discussion depuis les variables de route
	vars := mux.Vars(r)
	idStr := vars["id"]
	log.Printf("[DEBUG] GetThread - Thread ID depuis l'URL: %s", idStr)

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("[DEBUG] GetThread - Format d'ID invalide: %v", err)
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid thread ID",
		})
		return
	}

	// Récupérer le fil de discussion
	log.Printf("[DEBUG] GetThread - Tentative de récupération du thread ID: %d", id)
	thread, err := models.GetThread(c.DB, id)
	if err != nil {
		log.Printf("[DEBUG] GetThread - Erreur lors de la récupération du thread: %v", err)
		middleware.SendJSON(w, http.StatusNotFound, middleware.Response{
			Status:  "error",
			Message: "Thread not found",
		})
		return
	}

	log.Printf("[DEBUG] GetThread - Thread récupéré avec succès:")
	log.Printf("[DEBUG] GetThread - ID: %d", thread.ID)
	log.Printf("[DEBUG] GetThread - Titre: %s", thread.Title)
	log.Printf("[DEBUG] GetThread - Description: %s", thread.Description)
	log.Printf("[DEBUG] GetThread - Tags: %s", thread.Tags)
	log.Printf("[DEBUG] GetThread - AuthorID: %d", thread.AuthorID)
	log.Printf("[DEBUG] GetThread - MessageCount: %d", thread.MessageCount)

	if thread.Author != nil {
		log.Printf("[DEBUG] GetThread - Données de l'auteur:")
		log.Printf("[DEBUG] GetThread - Username: %s", thread.Author.Username)
		log.Printf("[DEBUG] GetThread - Email: %s", thread.Author.Email)
		log.Printf("[DEBUG] GetThread - Role: %s", thread.Author.Role)
	} else {
		log.Printf("[DEBUG] GetThread - ATTENTION: Author est nil!")
	}

	// Récupérer les messages
	log.Printf("[DEBUG] GetThread - Tentative de récupération des messages")
	messages, err := models.ListMessages(c.DB, id, 1, 10, "newest")
	if err != nil {
		log.Printf("[DEBUG] GetThread - Erreur lors de la récupération des messages: %v", err)
		messages = []*models.Message{} // Initialiser avec un slice vide en cas d'erreur
	} else {
		log.Printf("[DEBUG] GetThread - Nombre de messages récupérés: %d", len(messages))
		for i, msg := range messages {
			log.Printf("[DEBUG] GetThread - Message %d:", i+1)
			log.Printf("[DEBUG] GetThread - ID: %d", msg.ID)
			log.Printf("[DEBUG] GetThread - Content: %s", msg.Content)
			if msg.Author != nil {
				log.Printf("[DEBUG] GetThread - AuthorName: %s", msg.Author.Username)
			}
			log.Printf("[DEBUG] GetThread - CreatedAt: %s", msg.CreatedAt)
			log.Printf("[DEBUG] GetThread - Likes: %d", msg.Likes)
			log.Printf("[DEBUG] GetThread - Dislikes: %d", msg.Dislikes)
		}
	}

	// Formater les tags
	tags := strings.Split(thread.Tags, ",")
	if len(tags) == 1 && tags[0] == "" {
		tags = []string{}
	}
	log.Printf("[DEBUG] GetThread - Tags formatés: %v", tags)

	// Formater la date
	createdAt := thread.CreatedAt.Format("02/01/2006 à 15:04")
	log.Printf("[DEBUG] GetThread - Date formatée: %s", createdAt)

	// Préparer les données pour la réponse JSON
	responseData := map[string]interface{}{
		"id":          thread.ID,
		"title":       thread.Title,
		"description": thread.Description,
		"author":      thread.Author, // ou map[string]interface{} si besoin de filtrer
		"created_at":  thread.CreatedAt,
		"tags":        tags,
		"views":       thread.MessageCount,
		"messages":    messages,
	}

	log.Printf("[DEBUG] GetThread - Envoi de la réponse JSON")
	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   responseData,
	})
	return
}

// ListThreads récupère la liste des discussions
func (c *ThreadController) ListThreads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Récupérer les paramètres de pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 {
		perPage = 10
	}

	// Récupérer les paramètres de filtrage
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "open"
	}

	visibility := r.URL.Query().Get("visibility")
	if visibility == "" {
		visibility = "public"
	}

	// Récupérer les discussions
	threads, err := models.ListThreads(c.DB, page, perPage, status, visibility)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting threads",
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data: map[string]interface{}{
			"threads": threads,
		},
	})
}

// UpdateThread gère la mise à jour d'un fil de discussion
func (c *ThreadController) UpdateThread(w http.ResponseWriter, r *http.Request) {
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

	// Récupérer l'ID du fil de discussion
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid thread ID",
		})
		return
	}

	// Récupérer le fil de discussion
	thread, err := models.GetThread(c.DB, id)
	if err != nil {
		middleware.SendJSON(w, http.StatusNotFound, middleware.Response{
			Status:  "error",
			Message: "Thread not found",
		})
		return
	}

	// Vérifier les permissions
	if thread.AuthorID != claims.UserID && claims.Role != "admin" {
		middleware.SendJSON(w, http.StatusForbidden, middleware.Response{
			Status:  "error",
			Message: "Not authorized to update this thread",
		})
		return
	}

	var input struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
		Status      string   `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Mettre à jour le fil de discussion
	thread.Title = input.Title
	thread.Description = input.Description
	thread.Tags = strings.Join(input.Tags, ",")
	if input.Status != "" {
		thread.Status = string(models.ThreadStatus(input.Status))
	}

	if err := thread.UpdateThread(c.DB); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error updating thread: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   thread,
	})
}

// DeleteThread gère la suppression d'un fil de discussion
func (c *ThreadController) DeleteThread(w http.ResponseWriter, r *http.Request) {
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

	// Récupérer l'ID du fil de discussion
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid thread ID",
		})
		return
	}

	// Récupérer le fil de discussion
	thread, err := models.GetThread(c.DB, id)
	if err != nil {
		middleware.SendJSON(w, http.StatusNotFound, middleware.Response{
			Status:  "error",
			Message: "Thread not found",
		})
		return
	}

	// Vérifier les permissions
	if thread.AuthorID != claims.UserID && claims.Role != "admin" {
		middleware.SendJSON(w, http.StatusForbidden, middleware.Response{
			Status:  "error",
			Message: "Not authorized to delete this thread",
		})
		return
	}

	// Supprimer le fil de discussion
	if err := models.DeleteThread(c.DB, id); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error deleting thread: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status:  "success",
		Message: "Thread deleted successfully",
	})
}

// SearchThreads gère la recherche de fils de discussion
func (c *ThreadController) SearchThreads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Récupérer les paramètres de recherche
	query := r.URL.Query().Get("q")
	if query == "" {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Search query is required",
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

	// Rechercher les fils de discussion
	threads, err := models.SearchThreads(c.DB, query, page, perPage)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error searching threads: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   threads,
	})
}

// AdminUpdateThread permet à un admin de mettre à jour un fil de discussion
func (c *ThreadController) AdminUpdateThread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
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

	var input struct {
		Status     string `json:"status"`
		Visibility string `json:"visibility"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	if err := models.AdminUpdateThread(c.DB, threadID, input.Status, input.Visibility); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error updating thread",
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status:  "success",
		Message: "Thread updated successfully",
	})
}

// AdminDeleteMessage permet à un admin de supprimer un message
func (c *ThreadController) AdminDeleteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	messageID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid message ID",
		})
		return
	}

	if err := models.DeleteMessage(c.DB, messageID); err != nil {
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

// GetThreadMessages récupère tous les messages d'un fil de discussion
func (c *ThreadController) GetThreadMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Récupérer l'ID du fil de discussion depuis les variables de route
	vars := mux.Vars(r)
	idStr := vars["id"]
	log.Printf("[DEBUG] GetThreadMessages - Thread ID from URL: %s", idStr)

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("[DEBUG] GetThreadMessages - Invalid thread ID format: %v", err)
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
	if perPage < 1 {
		perPage = 10
	}

	// Récupérer le paramètre de tri
	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "newest"
	}

	// Récupérer les messages
	messages, err := models.ListMessages(c.DB, id, page, perPage, sortBy)
	if err != nil {
		log.Printf("[DEBUG] GetThreadMessages - Error getting messages: %v", err)
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error retrieving messages",
		})
		return
	}

	log.Printf("[DEBUG] GetThreadMessages - Messages found: %+v", messages)
	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data: map[string]interface{}{
			"messages": messages,
		},
	})
}

// CreateMessage crée un nouveau message dans un fil de discussion
func (c *ThreadController) CreateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Récupérer l'ID du fil de discussion depuis les variables de route
	vars := mux.Vars(r)
	threadIDStr := vars["id"]
	log.Printf("[DEBUG] CreateMessage - Thread ID from URL: %s", threadIDStr)

	threadID, err := strconv.ParseInt(threadIDStr, 10, 64)
	if err != nil {
		log.Printf("[DEBUG] CreateMessage - Invalid thread ID format: %v", err)
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid thread ID",
		})
		return
	}

	// Récupérer l'utilisateur depuis le token JWT
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	// Lire le contenu du message depuis le corps de la requête
	var request struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("[DEBUG] CreateMessage - Error decoding request body: %v", err)
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Créer le message
	message, err := models.CreateMessage(c.DB, threadID, claims.UserID, request.Content, "")
	if err != nil {
		log.Printf("[DEBUG] CreateMessage - Error creating message: %v", err)
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error creating message",
		})
		return
	}

	log.Printf("[DEBUG] CreateMessage - Message created: %+v", message)
	middleware.SendJSON(w, http.StatusCreated, middleware.Response{
		Status: "success",
		Data:   message,
	})
}

// LikeMessage ajoute un like à un message
func (c *ThreadController) LikeMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Récupérer l'ID du message depuis les variables de route
	vars := mux.Vars(r)
	messageIDStr := vars["id"]
	log.Printf("[DEBUG] LikeMessage - Message ID from URL: %s", messageIDStr)

	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		log.Printf("[DEBUG] LikeMessage - Invalid message ID format: %v", err)
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid message ID",
		})
		return
	}

	// Récupérer l'utilisateur depuis le token JWT
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	// Récupérer le message
	message, err := models.GetMessage(c.DB, messageID)
	if err != nil {
		log.Printf("[DEBUG] LikeMessage - Error getting message: %v", err)
		middleware.SendJSON(w, http.StatusNotFound, middleware.Response{
			Status:  "error",
			Message: "Message not found",
		})
		return
	}

	// Ajouter le like
	message.Likes++
	if err := message.UpdateMessage(c.DB); err != nil {
		log.Printf("[DEBUG] LikeMessage - Error liking message: %v", err)
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error liking message",
		})
		return
	}

	log.Printf("[DEBUG] LikeMessage - Message liked successfully")
	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status:  "success",
		Message: "Message liked successfully",
	})
}

// DislikeMessage ajoute un dislike à un message
func (c *ThreadController) DislikeMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Récupérer l'ID du message depuis les variables de route
	vars := mux.Vars(r)
	messageIDStr := vars["id"]
	log.Printf("[DEBUG] DislikeMessage - Message ID from URL: %s", messageIDStr)

	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		log.Printf("[DEBUG] DislikeMessage - Invalid message ID format: %v", err)
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid message ID",
		})
		return
	}

	// Récupérer l'utilisateur depuis le token JWT
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	// Récupérer le message
	message, err := models.GetMessage(c.DB, messageID)
	if err != nil {
		log.Printf("[DEBUG] DislikeMessage - Error getting message: %v", err)
		middleware.SendJSON(w, http.StatusNotFound, middleware.Response{
			Status:  "error",
			Message: "Message not found",
		})
		return
	}

	// Ajouter le dislike
	message.Dislikes++
	if err := message.UpdateMessage(c.DB); err != nil {
		log.Printf("[DEBUG] DislikeMessage - Error disliking message: %v", err)
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error disliking message",
		})
		return
	}

	log.Printf("[DEBUG] DislikeMessage - Message disliked successfully")
	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status:  "success",
		Message: "Message disliked successfully",
	})
}
