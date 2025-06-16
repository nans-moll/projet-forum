package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Fichiers statiques
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// REDIRECTION AUTOMATIQUE VERS REGISTER
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		// Redirection automatique vers la page register
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	})


	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {

			http.ServeFile(w, r, "./views/layouts/auth/register.html")
		} else if r.Method == "POST" {
			// Traitement du formulaire d'inscription
			fmt.Fprintf(w, "📝 Inscription reçue - Utilisateur créé avec succès !")
		}
	})

	// Page de connexion
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.ServeFile(w, r, "./views/layouts/auth/login.html")
		} else if r.Method == "POST" {
			fmt.Fprintf(w, "🔐 Connexion reçue")
		}
	})

	// Autres routes
	http.HandleFunc("/threads", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "💬 Page des discussions")
	})

	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "⚙️ Administration")
	})

	// Démarrage
	fmt.Println("🚀 CinéForum démarré sur http://localhost:8080")
	fmt.Println("📝 Redirection automatique vers la page d'inscription")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
