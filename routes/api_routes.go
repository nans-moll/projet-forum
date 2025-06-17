package routes

import (
	"net/http"

	"projet-forum/controllers"
	"projet-forum/middleware"

	"github.com/gorilla/mux"
)

// SetupAPIRoutes configure toutes les routes API
func SetupAPIRoutes(router *mux.Router, authController *controllers.AuthController, threadController *controllers.ThreadController, statsController *controllers.StatsController, userController *controllers.UserController) {
	// Routes publiques
	router.HandleFunc("/api/auth/register", authController.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authController.Login).Methods("POST")
	router.HandleFunc("/api/stats", statsController.GetStats).Methods("GET")
	router.HandleFunc("/api/threads", threadController.ListThreads).Methods("GET")
	router.HandleFunc("/api/threads/{id}", threadController.GetThread).Methods("GET")
	router.HandleFunc("/api/threads/{id}/messages", threadController.GetThreadMessages).Methods("GET")

	// Routes protégées
	router.Handle("/api/users/me", middleware.AuthMiddleware(http.HandlerFunc(userController.GetUserProfile))).Methods("GET")
	router.Handle("/api/users/me", middleware.AuthMiddleware(http.HandlerFunc(userController.UpdateUserProfile))).Methods("PUT")
	router.Handle("/api/users/stats", middleware.AuthMiddleware(http.HandlerFunc(userController.GetUserStats))).Methods("GET")
	router.Handle("/api/users/threads", middleware.AuthMiddleware(http.HandlerFunc(userController.GetUserThreads))).Methods("GET")
	router.Handle("/api/users/messages", middleware.AuthMiddleware(http.HandlerFunc(userController.GetUserMessages))).Methods("GET")

	// Routes des fils de discussion
	router.Handle("/api/threads", middleware.AuthMiddleware(http.HandlerFunc(threadController.CreateThread))).Methods("POST")
	router.Handle("/api/threads/{id}", middleware.AuthMiddleware(http.HandlerFunc(threadController.UpdateThread))).Methods("PUT")
	router.Handle("/api/threads/{id}", middleware.AuthMiddleware(http.HandlerFunc(threadController.DeleteThread))).Methods("DELETE")
	router.Handle("/api/threads/{id}/messages", middleware.AuthMiddleware(http.HandlerFunc(threadController.CreateMessage))).Methods("POST")
	router.Handle("/api/messages/{id}/like", middleware.AuthMiddleware(http.HandlerFunc(threadController.LikeMessage))).Methods("POST")
	router.Handle("/api/messages/{id}/dislike", middleware.AuthMiddleware(http.HandlerFunc(threadController.DislikeMessage))).Methods("POST")

	// Routes des pages HTML
	router.HandleFunc("/threads", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	}).Methods("GET")
	router.HandleFunc("/threads/{id}", threadController.ShowThreadPage).Methods("GET")
	router.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	}).Methods("GET")
	router.Handle("/profile", middleware.AuthMiddleware(http.HandlerFunc(userController.ShowProfilePage))).Methods("GET")
}
