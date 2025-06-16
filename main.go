package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"projet-forum/controllers"
	"projet-forum/middleware"
)

func main() {
	// Charger les variables d'environnement
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erreur lors du chargement du fichier .env")
	}

	// Configuration de la base de données
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Connexion à la base de données
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Vérifier la connexion
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Créer le routeur
	router := mux.NewRouter()

	// Servir les fichiers statiques
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Configuration des templates
	tmpl := template.Must(template.ParseGlob("views/layouts/*.html"))
	tmplAuth := template.Must(template.ParseGlob("views/layouts/auth/*.html"))

	// Route pour la page d'accueil
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	}).Methods("GET")

	// Route pour la page alternative
	router.HandleFunc("/page1", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "Page_1.html", nil)
	}).Methods("GET")

	// Routes d'authentification
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		tmplAuth.ExecuteTemplate(w, "login.html", nil)
	}).Methods("GET")

	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		tmplAuth.ExecuteTemplate(w, "register.html", nil)
	}).Methods("GET")

	// Middleware CORS
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Initialisation des contrôleurs
	authController := &controllers.AuthController{DB: db}
	userController := &controllers.UserController{DB: db}
	threadController := &controllers.ThreadController{DB: db}
	messageController := &controllers.MessageController{DB: db}
	searchController := &controllers.SearchController{DB: db}

	// Routes publiques
	router.HandleFunc("/api/auth/register", authController.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authController.Login).Methods("POST")

	// Routes protégées
	router.Handle("/api/users/me", middleware.AuthMiddleware(http.HandlerFunc(userController.GetUserProfile))).Methods("GET")
	router.Handle("/api/users/me", middleware.AuthMiddleware(http.HandlerFunc(userController.UpdateUserProfile))).Methods("PUT")
	router.Handle("/api/users/stats", middleware.AuthMiddleware(http.HandlerFunc(userController.GetUserStats))).Methods("GET")
	router.Handle("/api/users/threads", middleware.AuthMiddleware(http.HandlerFunc(userController.GetUserThreads))).Methods("GET")
	router.Handle("/api/users/messages", middleware.AuthMiddleware(http.HandlerFunc(userController.GetUserMessages))).Methods("GET")

	// Routes des fils de discussion
	router.HandleFunc("/api/threads", threadController.ListThreads).Methods("GET")
	router.Handle("/api/threads", middleware.AuthMiddleware(http.HandlerFunc(threadController.CreateThread))).Methods("POST")
	router.HandleFunc("/api/threads/{id}", threadController.GetThread).Methods("GET")
	router.Handle("/api/threads/{id}", middleware.AuthMiddleware(http.HandlerFunc(threadController.UpdateThread))).Methods("PUT")
	router.Handle("/api/threads/{id}", middleware.AuthMiddleware(http.HandlerFunc(threadController.DeleteThread))).Methods("DELETE")

	// Routes des messages
	router.HandleFunc("/api/threads/{id}/messages", messageController.ListMessages).Methods("GET")
	router.Handle("/api/threads/{id}/messages", middleware.AuthMiddleware(http.HandlerFunc(messageController.CreateMessage))).Methods("POST")
	router.Handle("/api/messages/{id}", middleware.AuthMiddleware(http.HandlerFunc(messageController.UpdateMessage))).Methods("PUT")
	router.Handle("/api/messages/{id}", middleware.AuthMiddleware(http.HandlerFunc(messageController.DeleteMessage))).Methods("DELETE")
	router.Handle("/api/messages/{id}/vote", middleware.AuthMiddleware(http.HandlerFunc(messageController.VoteMessage))).Methods("POST")

	// Routes de recherche
	router.HandleFunc("/api/search", searchController.SearchThreads).Methods("GET")

	// Récupérer le port depuis les variables d'environnement
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Démarrer le serveur
	fmt.Printf("🚀 Serveur démarré sur http://localhost:%s\n", port)
	fmt.Printf("📝 Documentation disponible sur http://localhost:%s/api/docs\n", port)
	fmt.Printf("💻 Interface utilisateur disponible sur http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
