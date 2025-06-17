package routes

import (
	"net/http"

	"projet-forum/controllers"
	"projet-forum/middleware"

	"github.com/gorilla/mux"
)

// SetupAuthRoutes configure les routes d'authentification
func SetupAuthRoutes(router *mux.Router, authController *controllers.AuthController) {
	// Routes publiques
	router.HandleFunc("/auth/register", authController.ShowRegisterForm).Methods("GET")
	router.HandleFunc("/auth/login", authController.ShowLoginForm).Methods("GET")

	// Routes protégées
	router.Handle("/auth/profile", middleware.AuthMiddleware(http.HandlerFunc(authController.ShowProfile))).Methods("GET")
	router.Handle("/auth/settings", middleware.AuthMiddleware(http.HandlerFunc(authController.ShowSettings))).Methods("GET")
}
