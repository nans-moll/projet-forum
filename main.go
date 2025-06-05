package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// TemplateData représente les données passées aux templates
type TemplateData struct {
	Title string
	Data  interface{}
}

// loadTemplates charge tous les templates HTML
func loadTemplates() *template.Template {
	tmpl := template.New("")

	// Charger tous les fichiers .html du dossier views
	err := filepath.Walk("views", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err = tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal("Erreur lors du chargement des templates:", err)
	}

	return tmpl
}

// renderTemplate est une fonction helper pour rendre les templates
func renderTemplate(w http.ResponseWriter, tmpl *template.Template, name string, data *TemplateData) {
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Charger les templates
	templates := loadTemplates()

	// Servir les fichiers statiques
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route principale
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := &TemplateData{
			Title: "Forum",
			Data:  nil,
		}
		renderTemplate(w, templates, "index.html", data)
	})

	// Configuration du serveur
	port := ":8080"
	log.Printf("Serveur démarré sur http://localhost%s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Erreur lors du démarrage du serveur:", err)
	}
}
