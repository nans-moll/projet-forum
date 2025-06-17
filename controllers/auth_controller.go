package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"projet-forum/middleware"
	"projet-forum/models"
)

type AuthController struct {
	DB *sql.DB
}

// Register gère l'inscription des utilisateurs
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
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

	// Validation du mot de passe
	if len(input.Password) < 12 {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Password must be at least 12 characters long",
		})
		return
	}

	hasUpper := strings.ContainsAny(input.Password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(input.Password)

	if !hasUpper || !hasSpecial {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Password must contain at least one uppercase letter and one special character",
		})
		return
	}

	// Créer l'utilisateur
	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Role:     "user",
		Banned:   false,
	}

	// Hasher le mot de passe
	if err := user.SetPassword(input.Password); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error hashing password: " + err.Error(),
		})
		return
	}

	// Enregistrer l'utilisateur
	query := `
		INSERT INTO users (
			username, 
			email, 
			password_hash, 
			role, 
			is_banned, 
			thread_count, 
			message_count, 
			last_connection,
			created_at
		) VALUES (?, ?, ?, ?, ?, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	result, err := c.DB.Exec(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.Banned,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "username") {
				middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
					Status:  "error",
					Message: "Ce nom d'utilisateur est déjà pris",
				})
			} else if strings.Contains(err.Error(), "email") {
				middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
					Status:  "error",
					Message: "Cette adresse email est déjà utilisée",
				})
			} else {
				middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
					Status:  "error",
					Message: "Cette information est déjà utilisée",
				})
			}
		} else {
			middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
				Status:  "error",
				Message: "Erreur lors de la création du compte: " + err.Error(),
			})
		}
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting user ID: " + err.Error(),
		})
		return
	}

	user.ID = id

	middleware.SendJSON(w, http.StatusCreated, middleware.Response{
		Status: "success",
		Data:   user,
	})
}

// Login gère la connexion d'un utilisateur
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
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

// ShowLoginForm affiche le formulaire de connexion
func (c *AuthController) ShowLoginForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/auth/login.html")
}

// ShowRegisterForm affiche le formulaire d'inscription
func (c *AuthController) ShowRegisterForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/auth/register.html")
}

// ShowProfile affiche le profil de l'utilisateur
func (c *AuthController) ShowProfile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/auth/profile.html")
}

// ShowSettings affiche les paramètres de l'utilisateur
func (c *AuthController) ShowSettings(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/auth/settings.html")
}

// GetUserProfile récupère le profil de l'utilisateur
func (c *AuthController) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	user, err := models.GetUserByID(c.DB, userID)
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

// UpdateUserProfile met à jour le profil de l'utilisateur
func (c *AuthController) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}
	user, err := models.GetUserByID(c.DB, userID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting user",
		})
		return
	}
	user.Username = input.Username
	user.Email = input.Email
	if err := user.Update(c.DB); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error updating user",
		})
		return
	}
	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   user,
	})
}

// GetUserStats récupère les statistiques de l'utilisateur
func (c *AuthController) GetUserStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	stats, err := models.GetUserStats(c.DB, userID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting user stats",
		})
		return
	}
	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   stats,
	})
}

// GetUserThreads récupère les fils de discussion de l'utilisateur
func (c *AuthController) GetUserThreads(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	threads, err := models.GetUserThreads(c.DB, userID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting user threads",
		})
		return
	}
	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   threads,
	})
}

// GetUserMessages récupère les messages de l'utilisateur
func (c *AuthController) GetUserMessages(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	messages, err := models.GetUserMessages(c.DB, userID)
	if err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error getting user messages",
		})
		return
	}
	middleware.SendJSON(w, http.StatusOK, middleware.Response{
		Status: "success",
		Data:   messages,
	})
}
