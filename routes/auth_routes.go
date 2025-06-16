package routes

import (
	"net/http"
	"projet-forum/controllers"
)

func RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login", controllers.Login)
	mux.HandleFunc("/register", controllers.Register)
}
