package controllers

import (
	"fmt"
	"net/http"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "📨 Messages - À implémenter")
}
