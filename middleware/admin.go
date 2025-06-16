package middleware

import (
	"net/http"
)

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Vérifier les droits admin
		next.ServeHTTP(w, r)
	}
}
