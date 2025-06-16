package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"projet-forum/middleware"
	"projet-forum/models"
	"strconv"
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

	// Authentifier l'utilisateur
	user, err := models.GetUserByEmail(c.DB, input.Email)
	if err != nil {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Invalid credentials",
		})
		return
	}

	if !user.CheckPassword(input.Password) {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Invalid credentials",
		})
		return
	}

	// Générer le token JWT
	token, err := middleware.GenerateToken(user)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error generating token",
		})
		return
	}

	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data: map[string]string{
			"token": token,
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

// GetUserProfile récupère le profil d'un utilisateur
func (c *UserController) GetUserProfile(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
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

// GetUserThreads récupère les fils de discussion d'un utilisateur
func (c *UserController) GetUserThreads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
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

	threads, err := models.GetThreadsByAuthorID(userID, limit, offset)
	if err != nil {
		http.Error(w, "Error fetching user threads", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(threads)
}

// GetUserMessages récupère les messages d'un utilisateur
func (c *UserController) GetUserMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
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

	messages, err := models.GetMessagesByAuthorID(userID, limit, offset)
	if err != nil {
		http.Error(w, "Error fetching user messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
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
