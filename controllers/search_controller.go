package controllers

import (
	"encoding/json"
	"net/http"
	"projet-forum/models"
	"strconv"
)

// SearchThreads recherche des fils de discussion par titre ou tag
func SearchThreads(w http.ResponseWriter, r *http.Request) {
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
		limit = 10 // Valeur par défaut
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	// Recherche par tag
	threadsByTag, err := models.GetThreadsByTag(query, limit, offset)
	if err != nil {
		http.Error(w, "Error searching threads by tag", http.StatusInternalServerError)
		return
	}

	// Recherche par titre
	threadsByTitle, err := models.GetThreadsByTitle(query, limit, offset)
	if err != nil {
		http.Error(w, "Error searching threads by title", http.StatusInternalServerError)
		return
	}

	// Combiner les résultats en évitant les doublons
	threadMap := make(map[int]*models.Thread)
	for _, thread := range threadsByTag {
		threadMap[thread.ID] = thread
	}
	for _, thread := range threadsByTitle {
		threadMap[thread.ID] = thread
	}

	// Convertir la map en slice
	var results []*models.Thread
	for _, thread := range threadMap {
		results = append(results, thread)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
