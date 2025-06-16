package middleware

import (
	"net/http"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Vérifier l'authentification
		next.ServeHTTP(w, r)
	}
}
