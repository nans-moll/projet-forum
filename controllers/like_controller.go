package controllers

import (
	"fmt"
	"net/http"
)

func LikeMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "👍 Like - À implémenter")
}
