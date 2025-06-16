package controllers

import (
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "./views/layouts/auth/login.html")
	} else if r.Method == "POST" {
		fmt.Fprintf(w, "🔐 Login POST - À implémenter")
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "📝 Register - À implémenter")
}
