package main

import (
	"/authentification/pages"
	"fmt"
	"log"
	"net/http"
	"threads/templates"
)

// Fonction principale qui initialise le serveur et définit les routes.
func main() {
	templates.InitTemplates() // Chargement des templates

	// Gestion des fichiers statiques (CSS, images, font)
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Définition des routes du site

	http.HandleFunc("/challenge1", pages.Defi1)
	http.HandleFunc("/challenge2", pages.Defi2)
	http.HandleFunc("/challenge3", pages.Challenge3)
	http.HandleFunc("/challenge4", pages.Challenge4)
	http.HandleFunc("/challenge5", pages.Challenge5)
	http.HandleFunc("/challenge6", pages.Challenge6)
	http.HandleFunc("/portfabio", pages.Portfabio)
	http.HandleFunc("/dashboard", pages.TableauDeBord)
	http.HandleFunc("/team", pages.Team)

	http.HandleFunc("/all-defis", pages.Alldefis)
	http.HandleFunc("/", pages.HomePage)

	// Démarrage du serveur sur le port 8080
	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
