package main

import (
	"fmt"
	"log"
	"net/http"
)

// Fonction principale qui initialise le serveur et définit les routes.
func main() {
	// Gestion des fichiers statiques (CSS, JS, images)
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Routes principales du forum cinéma
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/threads", ThreadsPage)
	http.HandleFunc("/thread/", ThreadViewPage)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/register", RegisterPage)
	http.HandleFunc("/admin", AdminDashboard)
	http.HandleFunc("/profile", ProfilePage)
	http.HandleFunc("/search", SearchPage)

	// Routes API
	http.HandleFunc("/api/auth/login", HandleLogin)
	http.HandleFunc("/api/auth/register", HandleRegister)
	http.HandleFunc("/api/threads", HandleThreads)
	http.HandleFunc("/api/messages", HandleMessages)
	http.HandleFunc("/api/likes", HandleLikes)

	// Démarrage du serveur sur le port 8080
	fmt.Println("🎬 CinéForum démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handlers temporaires (à remplacer par vos vrais controllers)
func HomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./views/layouts/index.html")
}

func ThreadsPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./views/layouts/threads/index.html")
}

func ThreadViewPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./views/layouts/threads/show.html")
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./views/layouts/auth/login.html")
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./views/layouts/auth/register.html")
}

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./views/layouts/admin/dashboard.html")
}

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./views/layouts/users/profile.html")
}

func SearchPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./views/layouts/search/index.html")
}

// Handlers API temporaires
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "success", "token": "fake_jwt_token"}`))
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "success", "message": "Compte créé avec succès"}`))
}

func HandleThreads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"threads": []}`))
}

func HandleMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"messages": []}`))
}

func HandleLikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "success"}`))
}