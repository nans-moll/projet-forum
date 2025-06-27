package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"projet-forum/middleware"
	"projet-forum/models"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type UserController struct {
	DB *sql.DB
}

// Register gère l'inscription d'un nouvel utilisateur
func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Créer l'utilisateur
	user, err := models.CreateUser(c.DB, input.Username, input.Email, input.Password, "user")
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error creating user: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusCreated, middleware.Response{
		Status: "success",
		Data:   user,
	})
}

// Login gère la connexion d'un utilisateur
func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Vérifier les identifiants
	user, err := models.AuthenticateUser(c.DB, input.Username, input.Password)
	if err != nil {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Invalid credentials",
		})
		return
	}

	// Générer le token JWT
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error generating token",
		})
		return
	}

	// Créer le cookie de session
	sessionCookie := &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24 heures
	}
	http.SetCookie(w, sessionCookie)

	// Mettre à jour la dernière connexion
	if err := models.UpdateLastConnection(c.DB, user.ID); err != nil {
		fmt.Printf("[DEBUG] Login - Erreur lors de la mise à jour de la dernière connexion: %v\n", err)
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data: map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
			},
		},
	})
}

// GetUser gère la récupération d'un utilisateur
func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Récupérer l'ID de l'utilisateur
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid user ID",
		})
		return
	}

	// Récupérer l'utilisateur
	user, err := models.GetUserByID(c.DB, id)
	if err != nil {
		middleware.SendJSON(w, http.StatusNotFound, middleware.Response{
			Status:  "error",
			Message: "User not found",
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   user,
	})
}

// UpdateUser gère la mise à jour d'un utilisateur
func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	// Récupérer l'ID de l'utilisateur
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid user ID",
		})
		return
	}

	// Vérifier les permissions
	if id != claims.UserID && claims.Role != "admin" {
		middleware.SendJSON(w, http.StatusForbidden, middleware.Response{
			Status:  "error",
			Message: "Not authorized to update this user",
		})
		return
	}

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Récupérer l'utilisateur
	user, err := models.GetUserByID(c.DB, id)
	if err != nil {
		middleware.SendJSON(w, http.StatusNotFound, middleware.Response{
			Status:  "error",
			Message: "User not found",
		})
		return
	}

	// Mettre à jour l'utilisateur
	user.Username = input.Username
	user.Email = input.Email
	if input.Password != "" {
		user.SetPassword(input.Password)
	}

	if err := user.UpdateUser(c.DB); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error updating user: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   user,
	})
}

// DeleteUser gère la suppression d'un utilisateur
func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	// Récupérer l'ID de l'utilisateur
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid user ID",
		})
		return
	}

	// Vérifier les permissions
	if id != claims.UserID && claims.Role != "admin" {
		middleware.SendJSON(w, http.StatusForbidden, middleware.Response{
			Status:  "error",
			Message: "Not authorized to delete this user",
		})
		return
	}

	// Supprimer l'utilisateur
	if err := models.DeleteUser(c.DB, id); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error deleting user: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status:  "success",
		Message: "User deleted successfully",
	})
}

// ListUsers gère la liste des utilisateurs
func (c *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	// Vérifier l'authentification
	claims := middleware.GetUserFromContext(r)
	if claims == nil || claims.Role != "admin" {
		middleware.SendJSON(w, http.StatusForbidden, middleware.Response{
			Status:  "error",
			Message: "Not authorized to list users",
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

	// Récupérer la liste des utilisateurs
	users, err := models.ListUsers(c.DB, page, perPage)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error listing users: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   users,
	})
}

// GetUserProfile récupère les informations du profil de l'utilisateur
func (c *UserController) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	// Récupérer les informations de l'utilisateur
	user, err := models.GetUserByID(c.DB, claims.UserID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting user profile",
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   user,
	})
}

// UpdateUserProfile met à jour le profil d'un utilisateur
func (c *UserController) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Username       string `json:"username"`
		Email          string `json:"email"`
		ProfilePicture string `json:"profile_picture"`
		Biography      string `json:"biography"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Mettre à jour les champs
	user, err := models.GetUserByID(c.DB, claims.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.ProfilePicture != "" {
		user.ProfilePicture = req.ProfilePicture
	}
	if req.Biography != "" {
		user.Biography = req.Biography
	}

	if err := user.UpdateUser(c.DB); err != nil {
		http.Error(w, "Error updating user profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetUserStats récupère les statistiques d'un utilisateur
func (c *UserController) GetUserStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByID(c.DB, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	stats := struct {
		MessageCount int    `json:"message_count"`
		ThreadCount  int    `json:"thread_count"`
		LastLogin    string `json:"last_login"`
	}{
		MessageCount: user.MessageCount,
		ThreadCount:  user.ThreadCount,
		LastLogin:    user.LastConnection.Format("2006-01-02 15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetUserThreads récupère les discussions créées par un utilisateur
func (c *UserController) GetUserThreads(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur depuis l'URL
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid user ID",
		})
		return
	}

	// Récupérer les discussions de l'utilisateur
	query := `
		SELECT t.id, t.title, t.content, t.created_at, 
			   (SELECT COUNT(*) FROM messages m WHERE m.thread_id = t.id) as message_count
		FROM threads t
		WHERE t.author_id = ?
		ORDER BY t.created_at DESC
		LIMIT 20
	`

	rows, err := c.DB.Query(query, userID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error fetching user threads",
		})
		return
	}
	defer rows.Close()

	var threads []map[string]interface{}
	for rows.Next() {
		var thread struct {
			ID           int64     `json:"id"`
			Title        string    `json:"title"`
			Content      string    `json:"content"`
			CreatedAt    time.Time `json:"created_at"`
			MessageCount int       `json:"message_count"`
		}

		err := rows.Scan(&thread.ID, &thread.Title, &thread.Content, &thread.CreatedAt, &thread.MessageCount)
		if err != nil {
			continue
		}

		threads = append(threads, map[string]interface{}{
			"id":            thread.ID,
			"title":         thread.Title,
			"content":       thread.Content,
			"created_at":    thread.CreatedAt,
			"message_count": thread.MessageCount,
		})
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   threads,
	})
}

// GetUserMessages récupère les messages postés par un utilisateur
func (c *UserController) GetUserMessages(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur depuis l'URL
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid user ID",
		})
		return
	}

	// Récupérer les messages de l'utilisateur
	query := `
		SELECT m.id, m.content, m.created_at, m.thread_id, t.title as thread_title
		FROM messages m
		JOIN threads t ON m.thread_id = t.id
		WHERE m.author_id = ?
		ORDER BY m.created_at DESC
		LIMIT 20
	`

	rows, err := c.DB.Query(query, userID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error fetching user messages",
		})
		return
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		var message struct {
			ID          int64     `json:"id"`
			Content     string    `json:"content"`
			CreatedAt   time.Time `json:"created_at"`
			ThreadID    int64     `json:"thread_id"`
			ThreadTitle string    `json:"thread_title"`
		}

		err := rows.Scan(&message.ID, &message.Content, &message.CreatedAt, &message.ThreadID, &message.ThreadTitle)
		if err != nil {
			continue
		}

		messages = append(messages, map[string]interface{}{
			"id":           message.ID,
			"content":      message.Content,
			"created_at":   message.CreatedAt,
			"thread_id":    message.ThreadID,
			"thread_title": message.ThreadTitle,
		})
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   messages,
	})
}

// UploadAvatar gère l'upload de l'avatar utilisateur
func (c *UserController) UploadAvatar(w http.ResponseWriter, r *http.Request) {
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

	// Analyser le formulaire multipart
	err := r.ParseMultipartForm(5 * 1024 * 1024) // 5MB max
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "File too large",
		})
		return
	}

	file, header, err := r.FormFile("profile_picture")
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "No file uploaded",
		})
		return
	}
	defer file.Close()

	// Valider le type de fichier
	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "File must be an image",
		})
		return
	}

	// Créer le dossier uploads s'il n'existe pas
	uploadsDir := "static/uploads/avatars"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error creating upload directory",
		})
		return
	}

	// Générer un nom de fichier unique
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("avatar_%d_%d%s", claims.UserID, time.Now().Unix(), ext)
	filepath := fmt.Sprintf("%s/%s", uploadsDir, filename)

	// Sauvegarder le fichier
	dst, err := os.Create(filepath)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error saving file",
		})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error copying file",
		})
		return
	}

	// Mettre à jour l'utilisateur avec le nouveau chemin de l'avatar
	profilePictureURL := fmt.Sprintf("/static/uploads/avatars/%s", filename)
	_, err = c.DB.Exec("UPDATE users SET profile_picture = ? WHERE id = ?", profilePictureURL, claims.UserID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error updating user profile",
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data: map[string]string{
			"profile_picture": profilePictureURL,
		},
	})
}

// GetCurrentUser récupère les informations de l'utilisateur connecté
func (c *UserController) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	user, err := models.GetUserByID(c.DB, claims.UserID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting user",
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   user,
	})
}

// UpdateCurrentUser met à jour les informations de l'utilisateur connecté
func (c *UserController) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Vérifier que l'utilisateur existe
	existingUser, err := models.GetUserByID(c.DB, claims.UserID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting user",
		})
		return
	}

	// Mettre à jour les champs autorisés
	existingUser.Username = user.Username
	existingUser.Email = user.Email

	// Mettre à jour l'utilisateur
	if err := existingUser.UpdateUser(c.DB); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error updating user",
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   existingUser,
	})
}

// UpdatePassword met à jour le mot de passe de l'utilisateur connecté
func (c *UserController) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Vérifier que l'utilisateur existe
	user, err := models.GetUserByID(c.DB, claims.UserID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting user",
		})
		return
	}

	// Vérifier le mot de passe actuel
	if !user.CheckPassword(req.CurrentPassword) {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid current password",
		})
		return
	}

	// Mettre à jour le mot de passe
	user.SetPassword(req.NewPassword)
	if err := user.UpdateUser(c.DB); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error updating password",
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status:  "success",
		Message: "Password updated successfully",
	})
}

// DeleteMessageReaction supprime une réaction
func (c *UserController) DeleteMessageReaction(w http.ResponseWriter, r *http.Request) {
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

	// Récupérer les paramètres de la requête
	messageIDStr := r.URL.Query().Get("message_id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid message ID",
		})
		return
	}

	// Supprimer la réaction
	err = models.DeleteMessageReaction(c.DB, messageID, claims.UserID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error deleting message reaction: " + err.Error(),
		})
		return
	}

	// Mettre à jour les compteurs
	err = models.UpdateMessageReactionCount(c.DB, messageID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error updating message reaction count: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status:  "success",
		Message: "Message reaction deleted successfully",
	})
}

// AddMessageReaction ajoute une réaction
func (c *UserController) AddMessageReaction(w http.ResponseWriter, r *http.Request) {
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

	// Récupérer les paramètres de la requête
	messageIDStr := r.URL.Query().Get("message_id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid message ID",
		})
		return
	}

	// Ajouter la réaction
	err = models.AddMessageReaction(c.DB, messageID, claims.UserID, "like")
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error adding message reaction: " + err.Error(),
		})
		return
	}

	// Mettre à jour les compteurs
	err = models.UpdateMessageReactionCount(c.DB, messageID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error updating message reaction count: " + err.Error(),
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status:  "success",
		Message: "Message reaction added successfully",
	})
}

// ShowProfilePage affiche la page de profil de l'utilisateur
func (c *UserController) ShowProfilePage(w http.ResponseWriter, r *http.Request) {
	// Servir la page de profil - l'authentification se fait côté client
	log.Printf("Affichage de la page de profil pour: %s", r.URL.Path)
	http.ServeFile(w, r, "templates/users/profile.html")
}

// ShowLoginPage affiche la page de connexion
func (c *UserController) ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/users/login.html")
}

// ShowRegisterPage affiche la page d'inscription
func (c *UserController) ShowRegisterPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/users/register.html")
}

// ShowAdminDashboard affiche le tableau de bord administrateur
func (c *UserController) ShowAdminDashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/admin/dashboard.html")
}

// BanUser bannit un utilisateur
func (c *UserController) BanUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
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

	if err := models.BanUser(c.DB, userID); err != nil {
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
