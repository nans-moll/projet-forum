package middleware

import (
	"context"
	"net/http"
	"projet-forum/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("votre_clé_secrète_jwt") // À remplacer par une vraie clé secrète

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

// AuthMiddleware vérifie le token JWT et ajoute les informations de l'utilisateur au contexte
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := bearerToken[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Vérifier si l'utilisateur existe toujours
		user, err := models.GetUserByID(claims.UserID)
		if err != nil || user.IsBanned {
			http.Error(w, "User not found or banned", http.StatusUnauthorized)
			return
		}

		// Ajouter les informations de l'utilisateur au contexte
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminMiddleware vérifie si l'utilisateur a le rôle admin
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(*models.User)
		if !ok || user.Role != "admin" {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext récupère l'utilisateur depuis le contexte
func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value("user").(*models.User)
	return user, ok
}
