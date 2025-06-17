package routes

import (
	"net/http"

	"projet-forum/controllers"
	"projet-forum/middleware"

	"github.com/gorilla/mux"
)

// SetupAdminRoutes configure les routes d'administration
func SetupAdminRoutes(router *mux.Router, adminController *controllers.AdminController) {
	// Routes protégées par le middleware Admin
	router.Handle("/admin/dashboard", middleware.AdminMiddleware(http.HandlerFunc(adminController.ShowDashboard))).Methods("GET")
	router.Handle("/admin/users", middleware.AdminMiddleware(http.HandlerFunc(adminController.ListUsers))).Methods("GET")
	router.Handle("/admin/users/{id}/ban", middleware.AdminMiddleware(http.HandlerFunc(adminController.BanUser))).Methods("POST")
	router.Handle("/admin/threads", middleware.AdminMiddleware(http.HandlerFunc(adminController.ListThreads))).Methods("GET")
	router.Handle("/admin/threads/{id}/status", middleware.AdminMiddleware(http.HandlerFunc(adminController.UpdateThreadStatus))).Methods("PUT")
	router.Handle("/admin/threads/{id}", middleware.AdminMiddleware(http.HandlerFunc(adminController.DeleteThread))).Methods("DELETE")
	router.Handle("/admin/messages/{id}", middleware.AdminMiddleware(http.HandlerFunc(adminController.DeleteMessage))).Methods("DELETE")
}
