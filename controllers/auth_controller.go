package controllers

import (
	"crypto/sha512"
	"encoding/json"
	"net/http"
	"projet-forum/middleware"
	"projet-forum/models"
	"regexp"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Identifier string `json:"identifier"` // Peut être username ou email
	Password   string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// Register gère l'inscription d'un nouvel utilisateur
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation du nom d'utilisateur
	if len(req.Username) < 3 {
		http.Error(w, "Username must be at least 3 characters long", http.StatusBadRequest)
		return
	}

	// Validation de l'email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Validation du mot de passe
	passwordRegex := regexp.MustCompile(`^(?=.*[A-Z])(?=.*[!@#$%^&*])(?=.*[0-9a-z]).{12,}$`)
	if !passwordRegex.MatchString(req.Password) {
		http.Error(w, "Password must be at least 12 characters long and contain at least one uppercase letter and one special character", http.StatusBadRequest)
		return
	}

	// Vérifier si l'utilisateur existe déjà
	if _, err := models.GetUserByEmail(req.Email); err == nil {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Hasher le mot de passe
	hasher := sha512.New()
	hasher.Write([]byte(req.Password))
	passwordHash := string(hasher.Sum(nil))

	// Créer le nouvel utilisateur
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}

	if err := user.CreateUser(); err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Générer le token JWT
	token, err := middleware.GenerateToken(user)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Retourner la réponse
	response := AuthResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Login gère la connexion d'un utilisateur
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hasher le mot de passe
	hasher := sha512.New()
	hasher.Write([]byte(req.Password))
	passwordHash := string(hasher.Sum(nil))

	// Chercher l'utilisateur par email ou username
	var user *models.User
	var err error

	// Essayer d'abord par email
	user, err = models.GetUserByEmail(req.Identifier)
	if err != nil {
		// Si pas trouvé par email, essayer par username
		user, err = models.GetUserByUsername(req.Identifier)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	}

	// Vérifier le mot de passe
	if user.PasswordHash != passwordHash {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Vérifier si l'utilisateur est banni
	if user.IsBanned {
		http.Error(w, "Account is banned", http.StatusForbidden)
		return
	}

	// Mettre à jour la dernière connexion
	if err := user.UpdateLastConnection(); err != nil {
		http.Error(w, "Error updating last connection", http.StatusInternalServerError)
		return
	}

	// Générer le token JWT
	token, err := middleware.GenerateToken(user)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Retourner la réponse
	response := AuthResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
