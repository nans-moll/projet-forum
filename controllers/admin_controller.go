package controllers

import (
	"fmt"
	"net/http"
)

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "⚙️ Dashboard Admin - À implémenter")
}
