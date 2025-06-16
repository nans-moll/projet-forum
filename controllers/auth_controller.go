package controllers

import (
	"database/sql"
	"encoding/json"
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
		user.Password,
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

// Login gère la connexion des utilisateurs
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendJSON(w, http.StatusMethodNotAllowed, middleware.Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	var input struct {
		Identifier string `json:"identifier"` // email ou username
		Password   string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.SendJSON(w, http.StatusBadRequest, middleware.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Essayer de trouver l'utilisateur par email ou username
	var user *models.User
	var err error

	if strings.Contains(input.Identifier, "@") {
		user, err = models.GetUserByEmail(c.DB, input.Identifier)
	} else {
		user, err = models.GetUserByUsername(c.DB, input.Identifier)
	}

	if err != nil {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Invalid credentials",
		})
		return
	}

	// Vérifier si l'utilisateur est banni
	if user.Banned {
		middleware.SendJSON(w, http.StatusForbidden, middleware.Response{
			Status:  "error",
			Message: "Account is banned",
		})
		return
	}

	// Vérifier le mot de passe
	if !user.CheckPassword(input.Password) {
		middleware.SendJSON(w, http.StatusUnauthorized, middleware.Response{
			Status:  "error",
			Message: "Invalid credentials",
		})
		return
	}

	// Mettre à jour la dernière connexion
	if err := user.UpdateLastConnection(c.DB); err != nil {
		middleware.SendJSON(w, http.StatusInternalServerError, middleware.Response{
			Status:  "error",
			Message: "Error updating last connection",
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
