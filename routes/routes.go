package routes

import (
	"net/http"
	"projet-forum/controllers"
	"projet-forum/middleware"
)

// SetupRoutes configure toutes les routes de l'application
func SetupRoutes(userController *controllers.UserController, threadController *controllers.ThreadController) {
	// Routes publiques
	http.HandleFunc("/", threadController.ShowThreadsPage)
	http.HandleFunc("/login", userController.ShowLoginPage)
	http.HandleFunc("/register", userController.ShowRegisterPage)
	http.HandleFunc("/api/login", userController.Login)
	http.HandleFunc("/api/register", userController.Register)
	http.HandleFunc("/api/threads", threadController.ListThreads)
	http.HandleFunc("/threads/show/", threadController.ShowThreadPage)

	// Routes protégées (nécessitent une authentification)
	http.Handle("/profile", middleware.AuthMiddleware(http.HandlerFunc(userController.ShowProfilePage)))
	http.Handle("/api/users/me", middleware.AuthMiddleware(http.HandlerFunc(userController.GetUserProfile)))
	http.Handle("/api/threads/create", middleware.AuthMiddleware(http.HandlerFunc(threadController.CreateThread)))
	http.Handle("/api/threads/", middleware.AuthMiddleware(http.HandlerFunc(threadController.UpdateThread)))
	http.Handle("/api/threads/", middleware.AuthMiddleware(http.HandlerFunc(threadController.DeleteThread)))
	http.Handle("/api/threads/search", middleware.AuthMiddleware(http.HandlerFunc(threadController.SearchThreads)))

	// Routes admin
	http.Handle("/admin", middleware.AdminMiddleware(http.HandlerFunc(userController.ShowAdminDashboard)))
	http.Handle("/api/admin/users/", middleware.AdminMiddleware(http.HandlerFunc(userController.BanUser)))
	http.Handle("/api/admin/threads/", middleware.AdminMiddleware(http.HandlerFunc(threadController.AdminUpdateThread)))
	http.Handle("/api/admin/messages/", middleware.AdminMiddleware(http.HandlerFunc(threadController.AdminDeleteMessage)))

	// Servir les fichiers statiques
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}
