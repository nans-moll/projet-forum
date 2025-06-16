package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"projet-forum/models"

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
func GenerateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
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
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			SendJSON(w, http.StatusUnauthorized, Response{
				Status:  "error",
				Message: "Authorization header is required",
			})
			return
		}

		// Extraire le token du header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
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

		if err != nil || !token.Valid {
			SendJSON(w, http.StatusUnauthorized, Response{
				Status:  "error",
				Message: "Invalid token",
			})
			return
		}

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
