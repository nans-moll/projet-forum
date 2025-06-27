package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"projet-forum/controllers"
	"projet-forum/routes"
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

	// Initialiser les contrôleurs
	authController := &controllers.AuthController{DB: db}
	userController := &controllers.UserController{DB: db}
	threadController := &controllers.ThreadController{DB: db}
	statsController := &controllers.StatsController{DB: db}
	adminController := &controllers.AdminController{DB: db}

	// Créer le routeur
	router := mux.NewRouter()

	// Configuration des types MIME
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			switch {
			case strings.HasSuffix(path, ".js"):
				w.Header().Set("Content-Type", "application/javascript")
			case strings.HasSuffix(path, ".css"):
				w.Header().Set("Content-Type", "text/css")
			case strings.HasSuffix(path, ".html"):
				w.Header().Set("Content-Type", "text/html")
			}
			next.ServeHTTP(w, r)
		})
	})

	// Servir les fichiers statiques
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Configuration des templates
	tmpl := template.New("").Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})

	// Charger tous les templates
	templateFiles := []string{
		"templates/index.html",
		"templates/threads/show.html",
		"templates/threads/edit.html",
		"templates/threads/create.html",
		"templates/users/profile.html",
	}

	for _, file := range templateFiles {
		_, err := tmpl.ParseFiles(file)
		if err != nil {
			log.Printf("Erreur lors du chargement du template %s: %v", file, err)
		}
	}

	// Configuration des routes
	routes.SetupAPIRoutes(router, authController, threadController, statsController, userController)
	routes.SetupAuthRoutes(router, authController)
	routes.SetupAdminRoutes(router, adminController)

	// Route pour la page d'accueil
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
			log.Printf("Erreur lors du rendu de index.html: %v", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		}
	}).Methods("GET")

	// Route pour la page de discussion
	router.HandleFunc("/threads/show/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		threadID := vars["id"]
		log.Printf("Affichage de la discussion %s", threadID)

		if err := tmpl.ExecuteTemplate(w, "show.html", nil); err != nil {
			log.Printf("Erreur lors du rendu de show.html: %v", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		}
	}).Methods("GET")

	// Route pour la page de liste des discussions
	router.HandleFunc("/threads", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Affichage de la page d'accueil")
		if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
			log.Printf("Erreur lors du rendu de index.html: %v", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		}
	}).Methods("GET")

	// Route de test pour l'accès au profil
	router.HandleFunc("/test-profile", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test-profile.html")
	}).Methods("GET")

	// Route de profil temporaire pour déboguer (sans middleware)
	router.HandleFunc("/profile-debug", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Accès à /profile-debug")
		http.ServeFile(w, r, "templates/users/profile.html")
	}).Methods("GET")

	// Démarrer le serveur
	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
