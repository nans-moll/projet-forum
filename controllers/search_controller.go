package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"projet-forum/models"
	"strconv"
)

type SearchController struct {
	DB *sql.DB
}

func (c *SearchController) SearchThreads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	// Rechercher les fils de discussion
	threads, err := models.SearchThreads(c.DB, query, limit, offset)
	if err != nil {
		http.Error(w, "Error searching threads", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(threads)
}
