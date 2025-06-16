package routes

import (
	"net/http"
	"projet-forum/controllers"
	"projet-forum/middleware"
)

func RegisterAdminRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin/dashboard", middleware.RequireAdmin(controllers.AdminDashboard))
}
