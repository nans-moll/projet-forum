package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"projet-forum/middleware"
	"projet-forum/models"
)

type ThreadController struct {
	DB *sql.DB
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

// GetThread gère la récupération d'un fil de discussion
func (c *ThreadController) GetThread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
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

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   thread,
	})
}

// ListThreads gère la liste des fils de discussion
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
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// Récupérer les filtres
	status := models.ThreadStatus(r.URL.Query().Get("status"))
	tag := r.URL.Query().Get("tag")

	// Récupérer la liste des fils de discussion
	threads, err := models.ListThreads(c.DB, page, perPage, string(status), tag)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error listing threads: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   threads,
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
