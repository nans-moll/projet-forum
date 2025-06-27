package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("votre_clé_secrète") // À remplacer par une vraie clé secrète

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken génère un token JWT pour un utilisateur
func GenerateToken(userID int64, username, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// SendJSON envoie une réponse JSON
func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// AuthMiddleware vérifie le token JWT
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[DEBUG] AuthMiddleware - URL: %s\n", r.URL.Path)
		fmt.Printf("[DEBUG] AuthMiddleware - Headers: %v\n", r.Header)

		// Pour les requêtes HTML (pages), on laisse passer et la vérification se fait côté client
		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			fmt.Printf("[DEBUG] AuthMiddleware - Requête HTML autorisée, vérification côté client\n")
			next.ServeHTTP(w, r)
			return
		}

		// Pour les requêtes API, vérifier le token JWT
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			fmt.Printf("[DEBUG] AuthMiddleware - Pas de header Authorization pour requête API\n")
			SendJSON(w, http.StatusUnauthorized, Response{
				Status:  "error",
				Message: "Authorization header is required",
			})
			return
		}

		// Extraire le token du header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Printf("[DEBUG] AuthMiddleware - Format de header invalide: %v\n", parts)
			SendJSON(w, http.StatusUnauthorized, Response{
				Status:  "error",
				Message: "Invalid authorization header format",
			})
			return
		}

		tokenString := parts[1]
		claims := &Claims{}

		// Vérifier le token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			fmt.Printf("[DEBUG] AuthMiddleware - Erreur de parsing du token: %v\n", err)
			SendJSON(w, http.StatusUnauthorized, Response{
				Status:  "error",
				Message: "Invalid token: " + err.Error(),
			})
			return
		}

		if !token.Valid {
			fmt.Printf("[DEBUG] AuthMiddleware - Token invalide\n")
			SendJSON(w, http.StatusUnauthorized, Response{
				Status:  "error",
				Message: "Invalid token",
			})
			return
		}

		fmt.Printf("[DEBUG] AuthMiddleware - Token valide pour l'utilisateur: %s (ID: %d)\n", claims.Username, claims.UserID)

		// Ajouter les claims au contexte
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminMiddleware vérifie si l'utilisateur est un administrateur
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserFromContext(r)
		if claims == nil || claims.Role != "admin" {
			SendJSON(w, http.StatusForbidden, Response{
				Status:  "error",
				Message: "Admin access required",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext récupère les claims de l'utilisateur depuis le contexte
func GetUserFromContext(r *http.Request) *Claims {
	claims, ok := r.Context().Value("user").(*Claims)
	if !ok {
		return nil
	}
	return claims
}
