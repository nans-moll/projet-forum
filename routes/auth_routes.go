package routes

import (
	"net/http"
	"projet-forum/controllers"
)

func SetupAuthRoutes(mux *http.ServeMux, authController *controllers.AuthController) {
	// Routes publiques
	mux.HandleFunc("/api/auth/register", authController.Register)
	mux.HandleFunc("/api/auth/login", authController.Login)
}
