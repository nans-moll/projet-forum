package main

import (
	"fmt"
	"log"
	"net/http"
	"projet-forum/config"
	"projet-forum/controllers"
	"projet-forum/middleware"
)

// Fonction principale qui initialise le serveur et définit les routes.
func main() {
	// Initialisation de la base de données
	if err := config.InitDB(); err != nil {
		log.Fatal("Erreur lors de l'initialisation de la base de données:", err)
	}

	// Gestion des fichiers statiques
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Routes d'authentification
	http.HandleFunc("/api/auth/register", controllers.Register)
	http.HandleFunc("/api/auth/login", controllers.Login)

	// Routes des fils de discussion
	http.Handle("/api/threads", middleware.AuthMiddleware(http.HandlerFunc(controllers.CreateThread)))
	http.HandleFunc("/api/threads/", controllers.GetThread)
	http.HandleFunc("/api/threads/tag", controllers.GetThreadsByTag)
	http.Handle("/api/threads/search", middleware.AuthMiddleware(http.HandlerFunc(controllers.SearchThreads)))
	http.Handle("/api/threads/update", middleware.AuthMiddleware(http.HandlerFunc(controllers.UpdateThread)))
	http.Handle("/api/threads/delete", middleware.AuthMiddleware(http.HandlerFunc(controllers.DeleteThread)))

	// Routes des messages
	http.Handle("/api/messages", middleware.AuthMiddleware(http.HandlerFunc(controllers.CreateMessage)))
	http.HandleFunc("/api/messages/", controllers.GetMessages)
	http.Handle("/api/messages/update", middleware.AuthMiddleware(http.HandlerFunc(controllers.UpdateMessage)))
	http.Handle("/api/messages/delete", middleware.AuthMiddleware(http.HandlerFunc(controllers.DeleteMessage)))
	http.Handle("/api/messages/like", middleware.AuthMiddleware(http.HandlerFunc(controllers.LikeMessage)))
	http.Handle("/api/messages/dislike", middleware.AuthMiddleware(http.HandlerFunc(controllers.DislikeMessage)))

	// Routes d'administration
	http.Handle("/api/admin/ban", middleware.AdminMiddleware(http.HandlerFunc(controllers.BanUser)))
	http.Handle("/api/admin/unban", middleware.AdminMiddleware(http.HandlerFunc(controllers.UnbanUser)))
	http.Handle("/api/admin/thread/status", middleware.AdminMiddleware(http.HandlerFunc(controllers.UpdateThreadStatus)))
	http.Handle("/api/admin/stats", middleware.AdminMiddleware(http.HandlerFunc(controllers.GetAdminStats)))

	// Démarrage du serveur
	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
