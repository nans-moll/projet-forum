package controllers

import (
	"fmt"
	"net/http"
)

func Search(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "🔍 Recherche - À implémenter")
}
