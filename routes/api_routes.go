package routes

import (
	"net/http"
	"projet-forum/controllers"
)

func RegisterAPIRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/threads", controllers.GetThreads)
	mux.HandleFunc("/api/messages", controllers.GetMessages)
	mux.HandleFunc("/api/search", controllers.Search)
}
