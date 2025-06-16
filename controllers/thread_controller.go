package controllers

import (
	"fmt"
	"net/http"
)

func GetThreads(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "💬 Liste des threads - À implémenter")
}

func CreateThread(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "➕ Créer un thread - À implémenter")
}
