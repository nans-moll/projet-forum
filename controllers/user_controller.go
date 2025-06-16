package controllers

import (
	"fmt"
	"net/http"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "👤 Profil utilisateur - À implémenter")
}
